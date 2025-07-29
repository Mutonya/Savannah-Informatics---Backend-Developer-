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

type AuthController struct {
	authService services.AuthService //Bussiness Logic Interface
}

// Constructor for AuthController
// initializes the controller with an `AuthService`.
// Dependency Injection through a constructor
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Getter for authService (primarily for testing this getter is for test access)
func (c *AuthController) AuthService() services.AuthService {
	return c.authService
}

// generateRandomState creates a base64-encoded random state string for OAuth2.
// Generates a random 16-byte sequence and encodes it in base64.
// This state is used in OAuth2 to prevent Cross Site Request Fogery attacks
func generateRandomState() string {
	b := make([]byte, 16)                       // Create 16-byte buffer
	_, _ = rand.Read(b)                         // Fill with cryptographically secure random bytes
	return base64.URLEncoding.EncodeToString(b) // Return URL-safe base64 encoded string
}

// Initiates the OAuth2 flow by redirecting the user to the OAuth provider's login page.
// generate the Base64 token  and store it into the cookie then get the url from the service and redirect
func (c *AuthController) Login(ctx *gin.Context) {
	// Generate secure random state
	state := generateRandomState()

	// Store state in a cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "oauthstate", // Cookie name
		Value:    state,
		HttpOnly: true,  // Prevent JavaScript access
		Secure:   false, // Set to true in production (HTTPS) Cross Site Scripting Security
		Path:     "/",   // Accessible to all paths
		MaxAge:   30000, // 5 minutes this is also part of security
	})
	// Get OAuth authorization URL from service
	authURL := c.authService.GetAuthCodeURL(state)
	log.Info().Str("state", state).Msg("Redirecting to OAuth provider")
	// Perform HTTP redirect
	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (c *AuthController) Callback(ctx *gin.Context) {
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

	// Exchange code for tokens and authenticate user
	// Authenticate user with service
	customer, accessToken, err := c.authService.Authenticate(ctx.Request.Context(), code)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authenticate user")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "authentication failed")
		return
	}

	log.Info().Str("email", customer.Email).Msg("User authenticated successfully")
	// Return success response with user data and token
	responses.SuccessResponse(ctx, http.StatusOK, gin.H{
		"customer":    customer,
		"accessToken": accessToken,
	})
}

// User Profile Handler
func (c *AuthController) Profile(ctx *gin.Context) {
	// Get customer ID from context (set by auth middleware)
	// the function returns (value and Boolean)
	customerID, exists := ctx.Get("customerID")
	if !exists {
		log.Warn().Msg("Unauthorized profile access attempt")
		responses.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}
	// Fetch customer details from service
	customer, err := c.authService.GetCustomerByID(customerID.(uint))
	if err != nil {
		log.Error().Err(err).Uint("customerID", customerID.(uint)).Msg("Failed to fetch customer profile")
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch profile")
		return
	}

	log.Info().Uint("customerID", customerID.(uint)).Msg("Profile fetched successfully")
	responses.SuccessResponse(ctx, http.StatusOK, customer)
}
