package services

import (
	"context"
	"strings"

	"gorm.io/gorm" //ORM for database operations (used for error handling).

	"github.com/Mutonya/Savanah/internal/config"
	"github.com/Mutonya/Savanah/internal/domain/models"
	"github.com/Mutonya/Savanah/internal/domain/repositories"
	"github.com/Mutonya/Savanah/pkg/oauth2"
)

// interface definition
// Defines the contract for authentication services  (method signatures)
type AuthService interface {
	GetAuthCodeURL(state string) string
	Authenticate(ctx context.Context, code string) (*models.Customer, string, error)
	ValidateToken(ctx context.Context, token string) (*models.Customer, error)
	GetCustomerByID(id uint) (*models.Customer, error)
}

// Encapsulation: Holds repository privately
// Service Implementation Struct
// combine values of different types into one logical unit
type authService struct {
	oauthProvider oauth2.OAuthProvider
	customerRepo  repositories.CustomerRepository
	config        *config.Config
}

// initialize  the service
func NewAuthService(oauthProvider oauth2.OAuthProvider, customerRepo repositories.CustomerRepository, config *config.Config) AuthService {
	return &authService{
		oauthProvider: oauthProvider, // handles OAuth2 flow
		customerRepo:  customerRepo,  // manages customer data
		config:        config,        //app config
	}
}

// Generates the OAuth2 authorization URL
// Delegates to OAuthProvider to construct the URL
// Returns URL to redirect user to identity provider
func (s *authService) GetAuthCodeURL(state string) string {
	return s.oauthProvider.GetAuthCodeURL(state)
}

// auth implimentation
func (s *authService) Authenticate(ctx context.Context, code string) (*models.Customer, string, error) {
	// Step 1: Exchange authorization code for tokens
	token, err := s.oauthProvider.Exchange(ctx, code)
	if err != nil {
		return nil, "", err
	}
	// Step 2: Verify ID token (JWT validation)
	idToken, err := s.oauthProvider.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, "", err
	}
	// Step 3: Extract claims from ID token
	var claims struct {
		Email   string `json:"email"` //user email {johndoe@gmail.com} google is our auth provider
		Name    string `json:"name"`  // user name { John Doe}
		Subject string `json:"sub"`   // Unique user ID from provider
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, "", err
	}
	// Step 4: Find or create customer
	customer, err := s.customerRepo.GetByOAuthID(claims.Subject)
	if err != nil {
		// handle  not found error
		if err == gorm.ErrRecordNotFound {
			// create customer
			customer = &models.Customer{
				OAuthID: claims.Subject,
				Email:   claims.Email,
			}
			//username {John Doe}
			// split the stringinto First and Last name
			names := strings.SplitN(claims.Name, " ", 2)
			if len(names) > 0 {

				customer.FirstName = names[0] // John
			}
			if len(names) > 1 {
				customer.LastName = names[1] // Doe
			}

			//Save to database
			if err := s.customerRepo.Create(customer); err != nil {
				return nil, "", err
			}
		} else {
			// other database Errors
			return nil, "", err
		}
	}

	//  return the customer Model and token {Authenticated User}
	// to meet Single responsibity you should use a mapper
	// the customer model should not be responsible for any other thing other than db access
	return customer, token.AccessToken, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*models.Customer, error) {
	// In a real implementation, we would validate the JWT token
	// For simplicity, we'll just get the customer by ID from the token claims
	// This would be replaced with proper JWT validation in production
	return s.customerRepo.GetByID(1) // {Place Holder}
}

// fetching the user with ID {Profile one  scenario }
func (s *authService) GetCustomerByID(id uint) (*models.Customer, error) {
	return s.customerRepo.GetByID(id)
}
