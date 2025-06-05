package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/rs/zerolog/log"
)

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

// ContentTypeMiddleware ensures JSON content type
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if contentType != "application/json" {
				log.Warn().Str("content-type", contentType).Msg("Invalid content type")
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
