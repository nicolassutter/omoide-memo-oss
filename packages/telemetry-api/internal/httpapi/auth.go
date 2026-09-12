package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const APIKeyHeader = "X-Api-Key"

func apiKeyMiddleware(expected string) gin.HandlerFunc {
	expectedB := []byte(expected)
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/v1/") {
			c.Next()
			return
		}

		got := []byte(c.GetHeader(APIKeyHeader))
		if len(expectedB) == 0 || subtle.ConstantTimeCompare(got, expectedB) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "missing or invalid api key"})
			return
		}
		c.Next()
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
