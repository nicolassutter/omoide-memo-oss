package httpapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"omoidememo.com/telemetry-api/internal/config"
	"omoidememo.com/telemetry-api/internal/telemetry"
)

const bodyLimit = 1 << 20

func NewRouter(cfg config.Config, store telemetry.Store) (http.Handler, huma.API) {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(corsMiddleware(cfg.AllowedOrigin, cfg.AllowCredentials))
	ginRouter.Use(bodyLimitMiddleware(bodyLimit))

	ginRouter.GET("/health", health)

	apiCfg := huma.DefaultConfig("Omoide Memo - Telemetry API", "1.0.0")
	apiCfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"PublicAPIKey": {
			Type:        "apiKey",
			In:          "header",
			Name:        PublicAPIKeyHeader,
			Description: "Public key shipped in mobile apps. Allows ingest only.",
		},
		"AdminAPIKey": {
			Type:        "apiKey",
			In:          "header",
			Name:        AdminAPIKeyHeader,
			Description: "Operator secret. Required for non-ingest endpoints.",
		},
	}
	api := humagin.New(ginRouter, apiCfg)

	publicKey := []byte(cfg.PublicAPIKey)
	adminKey := []byte(cfg.AdminAPIKey)

	ingestMiddlewares := huma.Middlewares{}
	fetchMiddlewares := huma.Middlewares{}
	if !cfg.DevMode {
		ingestMiddlewares = huma.Middlewares{createPublicAuthMiddleware(api, publicKey, adminKey)}
		fetchMiddlewares = huma.Middlewares{createAdminAuthMiddleware(api, adminKey)}
	}

	huma.Register(api, huma.Operation{
		OperationID:   "ingest-events",
		Method:        http.MethodPost,
		Path:          "/v1/events",
		Summary:       "Ingest a batch of telemetry events",
		Description:   "Accepts a telemetry batch envelope and persists each event.",
		DefaultStatus: http.StatusAccepted,
		Tags:          []string{"telemetry"},
		Middlewares:   ingestMiddlewares,
		Security: []map[string][]string{
			{"PublicAPIKey": {}},
			{"AdminAPIKey": {}},
		},
	}, createIngestHandler(store))

	huma.Register(api, huma.Operation{
		OperationID:   "fetch-events",
		Method:        http.MethodGet,
		Path:          "/v1/events",
		Summary:       "Fetch telemetry events",
		Description:   "Retrieves events filtered by a trackedAt time interval. Defaults to the last 30 days.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"telemetry"},
		Middlewares:   fetchMiddlewares,
		Security: []map[string][]string{
			{"AdminAPIKey": {}},
		},
	}, createFetchEventsHandler(store))

	return ginRouter, api
}

func corsMiddleware(allowedOrigin string, allowCredentials bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		if allowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, "+
				PublicAPIKeyHeader+", "+AdminAPIKeyHeader,
		)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
