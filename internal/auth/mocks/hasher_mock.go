package mocks

import "github.com/stretchr/testify/mock"

type HasherMock struct {
	mock.Mock
}

func (m *HasherMock) Hash(password string) string {
	args := m.Called(password)
	return args.String(0)
}

func (m *HasherMock) Verify(hash string, password string) bool {
	args := m.Called(hash, password)
	return args.Bool(0)
}
