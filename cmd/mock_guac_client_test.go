package cmd_test

import (
	"github.com/stretchr/testify/mock"
	guacTypes "github.com/techBeck03/guacamole-api-client/types"
)

type MockGuacClient struct {
	mock.Mock
}

func (m *MockGuacClient) CreateConnection(connection *guacTypes.GuacConnection) error {
	args := m.Called(connection)
	return args.Error(0)
}

func (m *MockGuacClient) CreateUser(user *guacTypes.GuacUser) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockGuacClient) DeleteUser(username string) error {
	args := m.Called(username)
	return args.Error(0)
}
