/*
Copyright © 2022 Jayson Grace <jayson.e.grace@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	log "github.com/cowdogmoo/guacinator/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	guac "github.com/techBeck03/guacamole-api-client"
	"github.com/techBeck03/guacamole-api-client/types"
)

// VncHost represents the parameters used to establish a new connection in Guacamole.
//
// **Attributes:**
// Name:     A string representing the name of the VNC host.
// IP:       A string representing the IP address of the VNC host.
// Port:     An integer representing the port to connect to on the VNC host.
// Password: A string representing the password for the VNC host.
type VncHost struct {
	Name     string
	IP       string
	Port     int
	Password string
}

// GuacService represents the interface for interacting with the Guacamole API.
//
// **Methods:**
// CreateGuacamoleConnection: Creates a new Guacamole connection.
// CreateAdminUser:           Creates a new Guacamole admin user.
// DeleteGuacUser:            Deletes a Guacamole user.
type GuacService interface {
	CreateGuacamoleConnection(vnchost VncHost) error
	CreateAdminUser(user, password string) error
	DeleteGuacUser(user string) error
}

// GuacServiceImpl represents the implementation of the GuacService interface.
type GuacServiceImpl struct{}

var (
	guacCfg     guac.Config
	guacClient  guac.Client
	guacAdminPW string
	user        string
	password    string
	vncHost     VncHost
	scheme      string
	guacURL     string

	// guacamoleCmd represents the guacamole command
	guacamoleCmd = &cobra.Command{
		Use:   "guacamole",
		Short: "Deploy and interface with Apache Guacamole.",

		Run: func(cmd *cobra.Command, args []string) {
			var err error
			/* Retrieve CLI args */
			guacURL, err = cmd.Flags().GetString("url")
			if err != nil {
				log.Error(
					"Failed to get Guacamole URL from CLI input: %v", err)
				cobra.CheckErr(err)
			}
			user, err = cmd.Flags().GetString("username")
			if err != nil {
				log.Error(
					"Failed to get username from CLI input: %v", err)
				cobra.CheckErr(err)
			}
			password, err = cmd.Flags().GetString("password")
			if err != nil {
				log.Error(
					"Failed to get password from CLI input: %v", err)
				cobra.CheckErr(err)
			}

			/* Unmarshal viper values */
			scheme = viper.GetString("guac.scheme")
			guacURL = viper.GetString("guac.url")
			vncHost.Port = viper.GetInt("guac.vnc_port")

			// Generate Guacamole connection config
			guacCfg = guac.Config{
				URL:                    fmt.Sprintf("%s://%s", scheme, guacURL),
				Username:               user,
				Password:               password,
				DisableTLSVerification: true,
			}
			guacService := &GuacServiceImpl{}

			/* Functionality provided by this cobra command. */

			// Get new guacadmin password from CLI (if applicable).
			guacAdminPW, err = cmd.Flags().GetString("guacadmin-pw")
			if err != nil {
				log.Error(
					"Failed to get input from CLI input: %v", err)
				cobra.CheckErr(err)
			}

			// Establish a client with Guacamole
			if err := connectGuac(guacCfg); err != nil {
				log.Error(err)
				cobra.CheckErr(err)
			}

			if guacAdminPW != "" {
				token, err := getToken()
				if err != nil {
					log.Error(err)
					cobra.CheckErr(err)
				}

				log.Info("Setting secure password for guacadmin")
				if err := setAdminPW(token, password, guacAdminPW); err != nil {
					log.Error(
						"Failed to set new Guacamole admin password: %v\n", err)
					cobra.CheckErr(err)
				}
				os.Exit(0)
			}

			vncHost.Name, err = cmd.Flags().GetString("connection")
			if err != nil {
				log.Error(
					"Failed to get input from CLI input: %v", err)
				cobra.CheckErr(err)
			}

			vncHost.Password, err = cmd.Flags().GetString("vnc-pw")
			if err != nil {
				log.Error(
					"Failed to get input from CLI input: %v", err)
				cobra.CheckErr(err)
			}

			vncHost.IP, err = cmd.Flags().GetString("vnc-ip")
			if err != nil {
				log.Error(
					"Failed to get input from CLI input: %v", err)
				cobra.CheckErr(err)
			}

			if vncHost.Name != "" && vncHost.Password != "" && vncHost.IP != "" {
				if err := guacService.CreateGuacamoleConnection(vncHost); err != nil {
					log.Error(
						"Failed to create %s connection in Guacamole: %v", vncHost.Name, err)
					cobra.CheckErr(err)
				}
				os.Exit(0)
			} else if vncHost.Name != "" && vncHost.Password == "" || vncHost.IP == "" {
				log.Error(
					"You must provide all required information to " +
						"add a new connection in Guacamole")
				cobra.CheckErr(err)
			}

			delUser, err := cmd.Flags().GetString("delete-user")
			if err != nil {
				log.Error(
					"Failed to get input from CLI input: %v", err)
				cobra.CheckErr(err)
			}
			if delUser != "" {
				if err := guacService.DeleteGuacUser(delUser); err != nil {
					log.Error(
						"Failed to delete %s from Guacamole: %v", delUser, err)
					cobra.CheckErr(err)
				}
				os.Exit(0)
			}

			newAdmin, err := cmd.Flags().GetString("new-admin")
			if err != nil {
				log.Error(
					"Failed to get input from CLI input: %v", err)
				cobra.CheckErr(err)
			}
			if newAdmin != "" {
				if err := guacService.CreateAdminUser(newAdmin, password); err != nil {
					log.Error(
						"Failed to create %s admin in Guacamole: %v", newAdmin, err)
					cobra.CheckErr(err)
				}
				os.Exit(0)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(guacamoleCmd)

	// Required inputs
	guacamoleCmd.Flags().StringP(
		"url", "l", "", "Guacamole URL.")
	if err := guacamoleCmd.MarkFlagRequired("url"); err != nil {
		log.Error(
			"Failed to mark required flag url: %v", err)
		cobra.CheckErr(err)
	}
	guacamoleCmd.Flags().StringP(
		"username", "u", "", "Username used to authenticate with Guacamole.")
	if err := guacamoleCmd.MarkFlagRequired("username"); err != nil {
		log.Error(
			"Failed to mark required flag username: %v", err)
		cobra.CheckErr(err)
	}
	guacamoleCmd.Flags().StringP(
		"password", "p", "", "Password used to authenticate with Guacamole.")
	if err := guacamoleCmd.MarkFlagRequired("password"); err != nil {
		log.Error(
			"Failed to mark required flag password: %v", err)
		cobra.CheckErr(err)
	}

	// Functionality provided by this cobra command.
	guacamoleCmd.Flags().StringP(
		"guacadmin-pw", "", "", "New password for the guacadmin user (we should not leave it as guacadmin).")
	guacamoleCmd.Flags().StringP(
		"connection", "", "", "Create a connection in Guacamole.")
	guacamoleCmd.Flags().StringP(
		"vnc-pw", "", "", "VNC password for device. Required to create a new connection.")
	guacamoleCmd.Flags().StringP(
		"vnc-ip", "", "", "IP address of host running VNC. Required to create a new connection.")
	guacamoleCmd.Flags().StringP(
		"delete-user", "", "", "Delete an input Guacamole user.")
	guacamoleCmd.Flags().StringP(
		"new-admin", "", "", "Create a new Guacamole admin user.")
}

// CreateGuacamoleConnection establishes a new connection
//
// in Guacamole using the provided VncHost information.
// CreatePackageDocs generates documentation for all Go packages in the current
// directory and its subdirectories. It traverses the file tree using a provided
// afero.Fs and Repo to create a new README.md file in each directory containing
// a Go package. It uses a specified template file for generating the README files.
//
// **Parameters:**
//
// vncHost: A VncHost struct containing the necessary information for the connection.
//
// **Returns:**
//
// error: An error if the connection cannot be created.
func (g *GuacServiceImpl) CreateGuacamoleConnection(vncHost VncHost) error {
	newConnection := types.GuacConnection{
		Name:             vncHost.Name,
		ParentIdentifier: "ROOT",
		Protocol:         "vnc",
		Attributes: types.GuacConnectionAttributes{
			MaxConnections:        "2",
			MaxConnectionsPerUser: "1",
		},
		Parameters: types.GuacConnectionParameters{
			Hostname: vncHost.IP,
			Port:     strconv.Itoa(vncHost.Port),
			Password: vncHost.Password,
		},
	}

	if err := guacClient.CreateConnection(&newConnection); err != nil {
		log.Error(
			"Failed to create %s connection in Guacamole: %v", vncHost.Name, err)
		return err
	}

	return nil
}

// CreateAdminUser creates a new admin user in
// Guacamole with the specified
// username and password.
//
// **Parameters:**
//
// user: A string representing the desired username for the new admin user.
//
// password: A string representing the desired password for the new admin user.
//
// **Returns:**
//
// error: An error if the admin user cannot be created.
func (g *GuacServiceImpl) CreateAdminUser(user, password string) error {
	newUser := types.GuacUser{
		Username: user,
		Password: password,
	}

	if err := guacClient.CreateUser(&newUser); err != nil {
		return err
	}

	permissionItems := []types.GuacPermissionItem{
		guacClient.NewAddSystemPermission(types.SystemPermissions{}.Administer()),
		guacClient.NewAddSystemPermission(types.SystemPermissions{}.CreateUser()),
		guacClient.NewAddSystemPermission(types.SystemPermissions{}.CreateConnection()),
		guacClient.NewAddSystemPermission(types.SystemPermissions{}.CreateConnectionGroup()),
		guacClient.NewAddSystemPermission(types.SystemPermissions{}.CreateSharingProfile()),
	}

	if err := guacClient.SetUserPermissions(newUser.Username, &permissionItems); err != nil {
		return err
	}

	fmt.Println("Successfully created admin user " + newUser.Username)

	return nil
}

// DeleteGuacUser removes a specified Guacamole user.
//
// **Parameters:**
//
// user: A string representing the username of the Guacamole user to be deleted.
//
// **Returns:**
//
// error: An error if the specified user cannot be deleted.
func (g *GuacServiceImpl) DeleteGuacUser(user string) error {
	if err := guacClient.DeleteUser(user); err != nil {
		return err
	}

	fmt.Println("Successfully created " + user)

	return nil

}

func getToken() (string, error) {
	var token string

	tokenPath := "api/tokens"
	resp, err := http.PostForm(fmt.Sprintf("%s://%s/%s", scheme, guacURL, tokenPath),
		url.Values{
			"username": []string{user},
			"password": []string{password},
		})

	if err != nil {
		log.Error(
			"Failed to get token from Guacamole: %v", err,
		)
		return token, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(
			"Failed to read response body from Guacamole: %v", err,
		)
		return token, err
	}

	var tokenresp types.AuthenticationResponse
	if err := json.Unmarshal(body, &tokenresp); err != nil {
		log.Error(
			"Failed to unmarshal response body from Guacamole: %v", err,
		)
		return token, err
	}
	token = tokenresp.AuthToken

	return token, nil
}

func connectGuac(cfg guac.Config) error {
	guacClient = guac.New(cfg)

	if err := guacClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Guacamole: %v", err)
	}

	return nil
}

func setAdminPW(token string, old string, new string) error {
	data := map[string]string{
		"oldPassword": old,
		"newPassword": new,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Error(
			"Failed to marshal payload for Guacamole: %v", err,
		)
		return err
	}

	adminPWResetURL := fmt.Sprintf("%s://%s/api/session/data/postgresql/users/guacadmin/password", scheme, guacURL)
	req, err := http.NewRequest("PUT", adminPWResetURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Error(
			"Failed to create request for Guacamole: %v", err,
		)
		return err
	}

	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("guacamole-token", token)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(
			"Failed to send request to Guacamole: %v", err,
		)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		log.Error(
			"Failed to change password: %v", err,
		)
	}

	return nil
}
