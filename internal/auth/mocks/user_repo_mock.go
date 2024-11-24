package mocks

import (
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) ExistsByLogin(login string) (bool, error) {
	args := m.Called(login)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepoMock) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepoMock) ExistsByUsername(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepoMock) FindByLogin(login string) (*model.User, error) {
	args := m.Called(login)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepoMock) Find(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepoMock) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepoMock) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepoMock) Add(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepoMock) Save(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepoMock) Delete(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
