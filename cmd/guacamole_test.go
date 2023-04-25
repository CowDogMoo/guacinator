package cmd_test

import (
	"testing"

	"github.com/cowdogmoo/guacinator/cmd"
	"github.com/techBeck03/guacamole-api-client/types"
)

var guacClient cmd.GuacClientInterface

type MockGuacClient struct{}

func (m *MockGuacClient) Connect() error {
	return nil
}

func (m *MockGuacClient) CreateConnection(connection *types.GuacConnection) error {
	return nil
}

func (m *MockGuacClient) CreateUser(user *types.GuacUser) error {
	return nil
}

func (m *MockGuacClient) SetUserPermissions(username string, permissionItems *[]types.GuacPermissionItem) error {
	return nil
}

func (m *MockGuacClient) DeleteUser(username string) error {
	return nil
}

func (m *MockGuacClient) NewAddSystemPermission(permission types.GuacPermissionType) types.GuacPermissionItem {
	return types.GuacPermissionItem{
		Op:    "add",
		Path:  "/permissions/system",
		Value: permission,
	}
}

func TestCreateGuacamoleConnection(t *testing.T) {
	mockClient := &MockGuacClient{}
	guacClient = mockClient

	vncHost := cmd.VncHost{Name: "Example", IP: "192.168.1.100", Port: 5900, Password: "password"}
	err := cmd.CreateGuacamoleConnection(guacClient, vncHost)

	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
}

func TestCreateAdminUser(t *testing.T) {
	mockClient := &MockGuacClient{}
	guacClient = mockClient

	username := "admin"
	password := "secure_password"
	err := cmd.CreateAdminUser(guacClient, username, password)

	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}
}

func TestDeleteGuacUser(t *testing.T) {
	mockClient := &MockGuacClient{}
	guacClient = mockClient

	username := "user_to_delete"
	err := cmd.DeleteGuacUser(guacClient, username)

	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
}
