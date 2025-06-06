package services_test

import (
	"context"
	"github.com/Mutonya/Savanah/internal/config"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Mutonya/Savanah/internal/domain/models"
	"github.com/Mutonya/Savanah/internal/domain/repositories/mocks"
	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/pkg/oauth2/mocks"
)

func TestAuthService_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	oauthMock := oauth2_mocks.NewMockOAuthProvider(ctrl)
	repoMock := repositories_mocks.NewMockCustomerRepository(ctrl)
	cfg := &config.Config{} // Add your config here

	// Create service with mocks
	authSvc := services.NewAuthService(oauthMock, repoMock, cfg)

	t.Run("successful authentication with new user", func(t *testing.T) {
		// Mock OAuth provider expectations
		oauthMock.EXPECT().
			Exchange(gomock.Any(), "valid_code").
			Return(&oauth2.Token{AccessToken: "test_token"}, nil)

		oauthMock.EXPECT().
			VerifyIDToken(gomock.Any(), gomock.Any()).
			Return(&oauth2.IDToken{}, nil)

		// Mock repository expectations
		repoMock.EXPECT().
			GetByOAuthID("new_user_id").
			Return(nil, gorm.ErrRecordNotFound)

		repoMock.EXPECT().
			Create(gomock.Any()).
			Return(nil)

		// Execute
		customer, token, err := authSvc.Authenticate(context.Background(), "valid_code")

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, "test_token", token)
		assert.Equal(t, "new_user_id", customer.OAuthID)
	})

	t.Run("authentication with existing user", func(t *testing.T) {
		// Test existing user scenario...
	})
}
