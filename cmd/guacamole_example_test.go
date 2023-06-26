package cmd_test

import (
	"fmt"
	"log"

	guacinator "github.com/cowdogmoo/guacinator/cmd"
)

func ExampleCreateGuacamoleConnection() {
	// Define a VncHost struct
	vncHost := guacinator.VncHost{Name: "Test", IP: "192.168.1.1", Port: 5900, Password: "test_password"}

	// Create an actual instance of Guacamole service
	// guacService := guacinator.NewGuacService()

	// Call the function with the struct
	// Here, instead of calling the actual function, we're just demonstrating how it would be used
	// In a real world context, replace this fmt.Println with guacService.CreateGuacamoleConnection(vncHost)
	fmt.Println("guacService.CreateGuacamoleConnection(vncHost)")
	log.Println("In actual use, replace the above fmt.Println with the actual method call using a real Guacamole service instance.")
	_ = vncHost
}

func ExampleCreateAdminUser() {
	// Define a username and password for the new admin user
	username := "test_admin"
	password := "test_password"

	// Create an actual instance of Guacamole service
	// guacService := guacinator.NewGuacService()

	// Call the function with the username and password
	// Here, instead of calling the actual function, we're just demonstrating how it would be used
	// In a real world context, replace this fmt.Println with guacService.CreateAdminUser(username, password)
	fmt.Println("guacService.CreateAdminUser(username, password)")
	log.Println("In actual use, replace the above fmt.Println with the actual method call using a real Guacamole service instance.")
	_ = username
	_ = password
}

func ExampleDeleteGuacUser() {
	// Define a username for the Guacamole user to be deleted
	username := "test_user"

	// Create an actual instance of Guacamole service
	// guacService := guacinator.NewGuacService()

	// Call the function with the username
	// Here, instead of calling the actual function, we're just demonstrating how it would be used
	// In a real world context, replace this fmt.Println with guacService.DeleteGuacUser(username)
	fmt.Println("guacService.DeleteGuacUser(username)")
	log.Println("In actual use, replace the above fmt.Println with the actual method call using a real Guacamole service instance.")
	_ = username
}
