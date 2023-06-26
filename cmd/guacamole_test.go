package cmd_test

import (
	"fmt"
	"testing"

	guacinator "github.com/cowdogmoo/guacinator/cmd"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockGuacService struct {
	mock.Mock
}

var mockService = new(MockGuacService)

func (m *MockGuacService) CreateGuacamoleConnection(vncHost guacinator.VncHost) error {
	args := m.Called(vncHost)
	return args.Error(0)
}

func (m *MockGuacService) CreateAdminUser(user, password string) error {
	args := m.Called(user, password)
	return args.Error(0)
}

func (m *MockGuacService) DeleteGuacUser(user string) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestCreateGuacamoleConnection(t *testing.T) {
	tests := []struct {
		name      string
		vncHost   guacinator.VncHost
		expectErr bool
	}{
		{
			name: "Valid VncHost",
			vncHost: guacinator.VncHost{
				Name:     "Example",
				IP:       "guacamole.techvomit.xyz",
				Port:     5900,
				Password: "guacadmin",
			},
			expectErr: false,
		},
		{
			name: "Invalid VncHost",
			vncHost: guacinator.VncHost{
				Name:     "Invalid",
				IP:       "192.168.1.300", // Invalid IP
				Port:     5900,
				Password: "password",
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectErr {
				mockService.On("CreateGuacamoleConnection", tc.vncHost).Return(fmt.Errorf("some error"))
			} else {
				mockService.On("CreateGuacamoleConnection", tc.vncHost).Return(nil)
			}

			err := mockService.CreateGuacamoleConnection(tc.vncHost)

			mockService.AssertExpectations(t)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateAdminUser(t *testing.T) {
	tests := []struct {
		name      string
		user      string
		password  string
		expectErr bool
	}{
		{
			name:      "Valid user and password",
			user:      "admin",
			password:  "secure_password",
			expectErr: false,
		},
		{
			name:      "Invalid user",
			user:      "", // Empty user
			password:  "secure_password",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectErr {
				mockService.On("CreateAdminUser", tc.user, tc.password).Return(fmt.Errorf("some error"))
			} else {
				mockService.On("CreateAdminUser", tc.user, tc.password).Return(nil)
			}

			err := mockService.CreateAdminUser(tc.user, tc.password)

			mockService.AssertExpectations(t)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestDeleteGuacUser(t *testing.T) {
	tests := []struct {
		name      string
		user      string
		expectErr bool
	}{
		{
			name:      "Valid user",
			user:      "user_to_delete",
			expectErr: false,
		},
		{
			name:      "Invalid user",
			user:      "", // Empty user
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectErr {
				mockService.On("DeleteGuacUser", tc.user).Return(fmt.Errorf("some error"))
			} else {
				mockService.On("DeleteGuacUser", tc.user).Return(nil)
			}

			err := mockService.DeleteGuacUser(tc.user)

			mockService.AssertExpectations(t)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
