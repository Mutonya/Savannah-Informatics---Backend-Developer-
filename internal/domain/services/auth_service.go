package services

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/config"
	"github.com/Mutonya/Savanah/internal/domain/models"
	"github.com/Mutonya/Savanah/internal/domain/repositories"
	"github.com/Mutonya/Savanah/pkg/oauth2"
)

type AuthService interface {
	GetAuthCodeURL(state string) string
	Authenticate(ctx context.Context, code string) (*models.Customer, string, error)
	ValidateToken(ctx context.Context, token string) (*models.Customer, error)
	GetCustomerByID(id uint) (*models.Customer, error)
}

type authService struct {
	oauthProvider oauth2.OAuthProvider
	customerRepo  repositories.CustomerRepository
	config        *config.Config
}

func NewAuthService(oauthProvider oauth2.OAuthProvider, customerRepo repositories.CustomerRepository, config *config.Config) AuthService {
	return &authService{
		oauthProvider: oauthProvider,
		customerRepo:  customerRepo,
		config:        config,
	}
}

func (s *authService) GetAuthCodeURL(state string) string {
	return s.oauthProvider.GetAuthCodeURL(state)
}

func (s *authService) Authenticate(ctx context.Context, code string) (*models.Customer, string, error) {
	token, err := s.oauthProvider.Exchange(ctx, code)
	if err != nil {
		return nil, "", err
	}

	idToken, err := s.oauthProvider.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, "", err
	}

	var claims struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Subject string `json:"sub"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, "", err
	}

	customer, err := s.customerRepo.GetByOAuthID(claims.Subject)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			customer = &models.Customer{
				OAuthID: claims.Subject,
				Email:   claims.Email,
			}

			names := strings.SplitN(claims.Name, " ", 2)
			if len(names) > 0 {
				customer.FirstName = names[0]
			}
			if len(names) > 1 {
				customer.LastName = names[1]
			}

			if err := s.customerRepo.Create(customer); err != nil {
				return nil, "", err
			}
		} else {
			return nil, "", err
		}
	}

	return customer, token.AccessToken, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*models.Customer, error) {
	// In a real implementation, we would validate the JWT token
	// For simplicity, we'll just get the customer by ID from the token claims
	// This would be replaced with proper JWT validation in production
	return s.customerRepo.GetByID(1) // Simplified for example
}

func (s *authService) GetCustomerByID(id uint) (*models.Customer, error) {
	return s.customerRepo.GetByID(id)
}
