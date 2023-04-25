package cmd_test

import (
	"errors"
	"testing"

	"github.com/cowdogmoo/guacinator/cmd"
)

func TestCreateGuacamoleConnection(t *testing.T) {
	tests := []struct {
		name     string
		vncHost  cmd.VncHost
		expected error
	}{
		{
			name: "Valid Connection",
			vncHost: cmd.VncHost{
				Name:     "Example",
				IP:       "192.168.1.100",
				Port:     5900,
				Password: "password",
			},
			expected: nil,
		},
		{
			name: "Invalid Connection",
			vncHost: cmd.VncHost{
				Name:     "Invalid",
				IP:       "192.168.1.999",
				Port:     0,
				Password: "",
			},
			expected: errors.New("Invalid Connection"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := cmd.CreateGuacamoleConnection(tc.vncHost)
			if err != tc.expected {
				t.Errorf("Expected error %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestCreateAdminUser(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		password string
		expected error
	}{
		{
			name:     "Valid Admin",
			user:     "admin",
			password: "secure_password",
			expected: nil,
		},
		{
			name:     "Invalid Admin",
			user:     "",
			password: "",
			expected: errors.New("Invalid User or Password"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := cmd.CreateAdminUser(tc.user, tc.password)
			if err != tc.expected {
				t.Errorf("Expected error %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestDeleteGuacUser(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		expected error
	}{
		{
			name:     "Valid User",
			user:     "user_to_delete",
			expected: nil,
		},
		{
			name:     "Invalid User",
			user:     "",
			expected: errors.New("Invalid User"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := cmd.DeleteGuacUser(tc.user)
			if err != tc.expected {
				t.Errorf("Expected error %v, got %v", tc.expected, err)
			}
		})
	}
}
