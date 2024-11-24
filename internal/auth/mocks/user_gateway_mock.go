package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UserGatewayMock struct {
	mock.Mock
}

func (m *UserGatewayMock) CreateUser(ctx context.Context, userID, role uint, username, email string) error {
	args := m.Called(ctx, userID, role, username, email)
	return args.Error(0)
}
