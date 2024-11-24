package mocks

import (
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"github.com/stretchr/testify/mock"
)

type UserActivationTokenServiceMock struct {
	mock.Mock
}

func (m *UserActivationTokenServiceMock) Generate(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *UserActivationTokenServiceMock) Get(userID uint) (*model.ActivationToken, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.ActivationToken), args.Error(0)
}

func (m *UserActivationTokenServiceMock) GetByToken(activationToken string) (*model.ActivationToken, error) {
	args := m.Called(activationToken)
	return args.Get(0).(*model.ActivationToken), args.Error(1)
}

func (m *UserActivationTokenServiceMock) Delete(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
