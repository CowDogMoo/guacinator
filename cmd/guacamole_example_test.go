package cmd_test

import (
	"fmt"

	guacinator "github.com/cowdogmoo/guacinator/cmd"
)

func ExampleGuacServiceImpl_CreateGuacamoleConnection() {
	// Define a VncHost struct
	vncHost := guacinator.VncHost{Name: "Test", IP: "192.168.1.1", Port: 5900, Password: "test_password"}

	// Create an actual instance of Guacamole service
	// guacService := guacinator.NewGuacService()

	// Call the function with the struct
	// Here, instead of calling the actual function, we're just demonstrating how it would be used
	fmt.Println("guacService.CreateGuacamoleConnection(vncHost)")

	_ = vncHost
	// Output:
	// guacService.CreateGuacamoleConnection(vncHost)
}

func ExampleGuacServiceImpl_CreateAdminUser() {
	// Define a username and password for the new admin user
	username := "test_admin"
	password := "test_password"

	// Create an actual instance of Guacamole service
	// guacService := guacinator.NewGuacService()

	// Call the function with the username and password
	// Here, instead of calling the actual function, we're just demonstrating how it would be used
	fmt.Println("guacService.CreateAdminUser(username, password)")

	_ = username
	_ = password
	// Output:
	// guacService.CreateAdminUser(username, password)
}

func ExampleGuacServiceImpl_DeleteGuacUser() {
	// Define a username for the Guacamole user to be deleted
	username := "test_user"

	// Create an actual instance of Guacamole service
	// guacService := guacinator.NewGuacService()

	// Call the function with the username
	// Here, instead of calling the actual function, we're just demonstrating how it would be used
	fmt.Println("guacService.DeleteGuacUser(username)")

	_ = username
	// Output:
	// guacService.DeleteGuacUser(username)
}
