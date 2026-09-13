package httpapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"omoidememo.com/telemetry-api/internal/config"
	"omoidememo.com/telemetry-api/internal/telemetry"
)

const (
	bodyLimit = 1 << 20
)

func NewRouter(cfg config.Config, store telemetry.Store) (http.Handler, huma.API) {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(corsMiddleware(cfg.AllowedOrigin))

	ginRouter.GET("/health", health)

	ginRouter.Use(apiKeyMiddleware(cfg.APIKey))
	ginRouter.Use(bodyLimitMiddleware(bodyLimit))

	api := humagin.New(ginRouter, huma.DefaultConfig("Omoide Memo - Telemetry API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID:   "ingest-events",
		Method:        http.MethodPost,
		Path:          "/v1/events",
		Summary:       "Ingest a batch of telemetry events",
		Description:   "Accepts a telemetry batch envelope and persists each event.",
		DefaultStatus: http.StatusAccepted,
		Tags:          []string{"telemetry"},
	}, createIngestHandler(store))

	huma.Register(api, huma.Operation{
		OperationID:   "fetch-events",
		Method:        http.MethodGet,
		Path:          "/v1/events",
		Summary:       "Fetch telemetry events",
		Description:   "Retrieves events filtered by a trackedAt time interval. Defaults to the last 30 days.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"telemetry"},
	}, createFetchEventsHandler(store))

	return ginRouter, api
}

func corsMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Api-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
