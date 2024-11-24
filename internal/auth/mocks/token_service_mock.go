package mocks

import "github.com/stretchr/testify/mock"

type TokenServiceMock struct {
	mock.Mock
}

func (m *TokenServiceMock) Generate(userID uint) (string, string, error) {
	args := m.Called(userID)
	return args.String(0), args.String(1), args.Error(2)
}

// Add other mock methods for UserActivationTokenServiceMock, HasherMock, and UserGatewayMock as needed
