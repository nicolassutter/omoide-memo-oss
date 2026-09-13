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
