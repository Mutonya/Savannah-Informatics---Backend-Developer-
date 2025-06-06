package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/internal/utils/responses"
)

type TestAuthController struct {
	authService services.AuthService
}

func NewAuthControllerTest(authService services.AuthService) *TestAuthController {
	return &TestAuthController{authService: authService}
}

// AuthService getter
func (c *TestAuthController) AuthService() services.AuthService {
	return c.authService
}

// generateTestRandomState creates a base64-encoded random state string for OAuth2.
func generateTestRandomState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Login @Summary Login with OAuth2
// @Description Redirects to OAuth provider's login page
// @Tags auth
// @Accept  json
// @Produce  json
// @Param state query string true "State parameter for CSRF protection"
// @Success 302 {string} string "Redirect to OAuth provider"
// @Failure 400 {object} responses.ErrorResponse
// @Router /auth/login [get]
func (c *TestAuthController) Login(ctx *gin.Context) {
	// Generate secure random state
	state := generateTestRandomState()

	// Store state in a cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		HttpOnly: true,
		Secure:   false, // Set to true in production (HTTPS)
		Path:     "/",
		MaxAge:   300, // 5 minutes
	})

	authURL := c.authService.GetAuthCodeURL(state)
	log.Info().Str("state", state).Msg("Redirecting to OAuth provider")
	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback @Summary OAuth2 Callback
// @Description Handles OAuth2 callback from provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param code query string true "Authorization code from provider"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/callback [get]
func (c *TestAuthController) Callback(ctx *gin.Context) {
	// Validate state parameter
	stateFromQuery := ctx.Query("state")
	if stateFromQuery == "" {
		log.Warn().Msg("Callback attempt without state parameter")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "state parameter is required")
		return
	}

	// Get state from cookie
	cookie, err := ctx.Request.Cookie("oauthstate")
	if err != nil || cookie.Value != stateFromQuery {
		log.Warn().
			Str("expected", func() string {
				if cookie != nil {
					return cookie.Value
				}
				return ""
			}()).
			Str("received", stateFromQuery).
			Msg("Invalid or missing state")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "invalid state parameter")
		return
	}

	// Validate code
	code := ctx.Query("code")
	if code == "" {
		log.Warn().Msg("Callback attempt without code parameter")
		responses.ErrorResponse(ctx, http.StatusBadRequest, "code parameter is required")
		return
	}

	// Authenticate user with service
	customer, accessToken, err := c.authService.Authenticate(ctx.Request.Context(), code)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authenticate user")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "authentication failed")
		return
	}

	log.Info().Str("email", customer.Email).Msg("User authenticated successfully")
	responses.SuccessResponse(ctx, http.StatusOK, gin.H{
		"customer":    customer,
		"accessToken": accessToken,
	})
}

// Profile @Summary Get User Profile
// @Description Get authenticated user's profile
// @Tags auth
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.SuccessResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/profile [get]
func (c *TestAuthController) Profile(ctx *gin.Context) {
	customerID, exists := ctx.Get("customerID")
	if !exists {
		log.Warn().Msg("Unauthorized profile access attempt")
		responses.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	customer, err := c.authService.GetCustomerByID(customerID.(uint))
	if err != nil {
		log.Error().Err(err).Uint("customerID", customerID.(uint)).Msg("Failed to fetch customer profile")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch profile")
		return
	}

	log.Info().Uint("customerID", customerID.(uint)).Msg("Profile fetched successfully")
	responses.SuccessResponse(ctx, http.StatusOK, customer)
}
