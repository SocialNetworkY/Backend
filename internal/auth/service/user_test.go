package service

import (
	"errors"
	"testing"

	"github.com/SocialNetworkY/Backend/internal/auth/mocks"
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Login(t *testing.T) {
	userRepo := new(mocks.UserRepoMock)
	tokenService := new(mocks.TokenServiceMock)
	hasher := new(mocks.HasherMock)

	userID := uint(1)
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// Setup successful login scenario
	hashedPassword := "hashed_password123" // Hash of the password
	expectedAccessToken := "access_token"
	expectedRefreshToken := "refresh_token"

	userRepo.On("FindByLogin", username).Return(&model.User{
		ID:          userID,
		Email:       email,
		Username:    username,
		Password:    hashedPassword,
		IsActivated: true,
	}, nil)
	tokenService.On("Generate", userID).Return(expectedAccessToken, expectedRefreshToken, nil)
	hasher.On("Verify", hashedPassword, password).Return(true)

	service := NewUserService(userRepo, tokenService, new(mocks.UserActivationTokenServiceMock), hasher, new(mocks.UserGatewayMock))

	accessToken, refreshToken, err := service.Login(username, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccessToken, accessToken)
	assert.Equal(t, expectedRefreshToken, refreshToken)

	// Setup failure scenarios
	testCases := []struct {
		description string
		username    string
		password    string
		expectedErr error
	}{
		{
			description: "Invalid password",
			username:    username,
			password:    "wrongpassword123",
			expectedErr: ErrInvalidPassword,
		},
		{
			description: "User not activated",
			username:    username,
			password:    password,
			expectedErr: ErrUserNotActivated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			userRepo := new(mocks.UserRepoMock)
			hasher := new(mocks.HasherMock)
			service.repo = userRepo
			service.hasher = hasher

			userRepo.On("FindByLogin", tc.username).Return(&model.User{
				ID:          userID,
				Email:       email,
				Username:    tc.username,
				Password:    hashedPassword,
				IsActivated: !errors.Is(ErrUserNotActivated, tc.expectedErr),
			}, nil)
			hasher.On("Verify", hashedPassword, tc.password).Return(!errors.Is(ErrInvalidPassword, tc.expectedErr))

			accessToken, refreshToken, err := service.Login(tc.username, tc.password)
			assert.Error(t, err)
			assert.Empty(t, accessToken)
			assert.Empty(t, refreshToken)
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}
