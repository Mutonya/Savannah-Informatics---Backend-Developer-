package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/internal/utils/errors"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAPIError(http.StatusUnauthorized, "authorization header is required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAPIError(http.StatusUnauthorized, "invalid authorization header format"))
			return
		}

		token := parts[1]
		customer, err := authService.ValidateToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAPIError(http.StatusUnauthorized, "invalid token"))
			return
		}

		ctx.Set("customerID", customer.ID)
		ctx.Next()
	}
}
