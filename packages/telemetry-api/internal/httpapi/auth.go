package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

const (
	PublicAPIKeyHeader = "X-Api-Public-Key"
	AdminAPIKeyHeader  = "X-Admin-Api-Key"
)

func writeUnauthorized(api huma.API, ctx huma.Context) {
	_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing or invalid api key")
}

func createPublicAuthMiddleware(api huma.API, publicKey, adminKey []byte) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		presentedPublicKey := []byte(ctx.Header(PublicAPIKeyHeader))
		presentedAdminKey := []byte(ctx.Header(AdminAPIKeyHeader))
		publicMatchesPublic := subtle.ConstantTimeCompare(presentedPublicKey, publicKey) == 1
		publicMatchesAdmin := subtle.ConstantTimeCompare(presentedPublicKey, adminKey) == 1
		adminMatchesAdmin := subtle.ConstantTimeCompare(presentedAdminKey, adminKey) == 1

		if !publicMatchesPublic && !publicMatchesAdmin && !adminMatchesAdmin {
			writeUnauthorized(api, ctx)
			return
		}
		next(ctx)
	}
}

func createAdminAuthMiddleware(api huma.API, adminKey []byte) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		presentedAdminKey := []byte(ctx.Header(AdminAPIKeyHeader))
		if subtle.ConstantTimeCompare(presentedAdminKey, adminKey) != 1 {
			writeUnauthorized(api, ctx)
			return
		}
		next(ctx)
	}
}

func bodyLimitMiddleware(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/v1/") {
			c.Next()
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max)
		c.Next()
	}
}
