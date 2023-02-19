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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	guac "github.com/techBeck03/guacamole-api-client"
	"github.com/techBeck03/guacamole-api-client/types"
)

// VncHost is used to hold the information
// necessary to incorporate a host running VNC
// into Guacamole.
type VncHost struct {
	Name     string
	Password string
	IP       string
	Port     int
}

var (
	cfg         guac.Config
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
			/* Retrieve CLI args */
			guacURL, err = cmd.Flags().GetString("url")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get Guacamole URL from CLI input: %v", err)
				os.Exit(1)
			}
			user, err = cmd.Flags().GetString("username")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get username from CLI input: %v", err)
				os.Exit(1)
			}
			password, err = cmd.Flags().GetString("password")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get password from CLI input: %v", err)
				os.Exit(1)
			}

			/* Unmarshal viper values */
			scheme = viper.GetString("guac.scheme")
			guacURL = viper.GetString("guac.url")
			vncHost.Port = viper.GetInt("guac.vnc_port")

			// Generate Guacamole connection config
			cfg = guac.Config{
				URL:                    fmt.Sprintf("%s://%s", scheme, guacURL),
				Username:               user,
				Password:               password,
				DisableTLSVerification: true,
			}

			/* Functionality provided by this cobra command. */

			// Get new guacadmin password from CLI (if applicable).
			guacAdminPW, err = cmd.Flags().GetString("guacadmin-pw")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get input from CLI input: %v", err)
				os.Exit(1)
			}

			// Establish a client with Guacamole
			if err := connectGuac(cfg); err != nil {
				log.WithError(err).Errorf(
					"failed to connect to Guacamole: %v", err)
				os.Exit(1)
			}

			if guacAdminPW != "" {
				token, err := getToken()
				if err != nil {
					log.WithError(err).Errorf(
						"failed to retrieve token from Guacamole: %v", err)
					os.Exit(1)
				}

				log.Info("Setting secure password for guacadmin")
				if err := setAdminPW(token, password, guacAdminPW); err != nil {
					log.WithError(err).Errorf(
						"failed to install Guacamole: %v", err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			vncHost.Name, err = cmd.Flags().GetString("connection")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get input from CLI input: %v", err)
				os.Exit(1)
			}

			vncHost.Password, err = cmd.Flags().GetString("vnc-pw")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get input from CLI input: %v", err)
				os.Exit(1)
			}

			vncHost.IP, err = cmd.Flags().GetString("vnc-ip")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get input from CLI input: %v", err)
				os.Exit(1)
			}

			if vncHost.Name != "" && vncHost.Password != "" && vncHost.IP != "" {
				if err := CreateGuacamoleConnection(vncHost); err != nil {
					log.WithError(err).Errorf(
						"failed to create %s connection in Guacamole: %v", vncHost.Name, err)
					os.Exit(1)
				}
				os.Exit(0)
			} else if vncHost.Name != "" && vncHost.Password == "" || vncHost.IP == "" {
				log.Error(
					"you must provide all required information to " +
						"add a new connection in Guacamole")
				os.Exit(1)

			}

			delUser, err := cmd.Flags().GetString("delete-user")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get input from CLI input: %v", err)
				os.Exit(1)
			}
			if delUser != "" {
				if err := DeleteGuacUser(delUser); err != nil {
					log.WithError(err).Errorf(
						"failed to delete %s from Guacamole: %v", delUser, err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			newAdmin, err := cmd.Flags().GetString("new-admin")
			if err != nil {
				log.WithError(err).Errorf(
					"failed to get input from CLI input: %v", err)
				os.Exit(1)
			}
			if newAdmin != "" {
				if err := CreateAdminUser(newAdmin, password); err != nil {
					log.WithError(err).Errorf(
						"failed to create %s admin in Guacamole: %v", newAdmin, err)
					os.Exit(1)
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
	guacamoleCmd.MarkFlagRequired("url")
	guacamoleCmd.Flags().StringP(
		"username", "u", "", "Username used to authenticate with Guacamole.")
	guacamoleCmd.MarkFlagRequired("username")
	guacamoleCmd.Flags().StringP(
		"password", "p", "", "Password used to authenticate with Guacamole.")
	guacamoleCmd.MarkFlagRequired("password")

	// Functionality provided by this cobra command.
	guacamoleCmd.Flags().BoolP(
		"install", "", false, "Install Guacamole on k8s.")
	guacamoleCmd.Flags().BoolP(
		"destroy", "", false, "Destroy k8s Guacamole deployment.")
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

func getToken() (string, error) {
	var token string

	tokenPath := "api/tokens"
	resp, err := http.PostForm(fmt.Sprintf("%s://%s/%s", scheme, guacURL, tokenPath),
		url.Values{
			"username": []string{user},
			"password": []string{password},
		})

	if err != nil {
		log.WithError(err)
		return token, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err)
		return token, err
	}

	var tokenresp types.AuthenticationResponse
	if err := json.Unmarshal(body, &tokenresp); err != nil {
		log.WithError(err)
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
		log.WithError(err)
		return err
	}

	adminPWResetURL := fmt.Sprintf("%s://%s/api/session/data/postgresql/users/guacadmin/password", scheme, guacURL)
	req, err := http.NewRequest("PUT", adminPWResetURL, bytes.NewBuffer(payload))
	if err != nil {
		log.WithError(err)
		return err
	}

	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("guacamole-token", token)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		log.WithError(fmt.Errorf("failed to change password: %v", err))
	}

	return nil
}

// CreateGuacamoleConnection creates a connection in guacamole.
func CreateGuacamoleConnection(vncHost VncHost) error {

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
		return err
	}

	return nil
}

// CreateAdminUser creates an admin user with the input password.
func CreateAdminUser(user string, password string) error {
	newUser := types.GuacUser{
		Username: user,
		Password: password,
	}

	if err = guacClient.CreateUser(&newUser); err != nil {
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

// DeleteGuacUser deletes a Guacamole user.
func DeleteGuacUser(user string) error {

	if err = guacClient.DeleteUser(user); err != nil {
		return err
	}

	fmt.Println("Successfully created " + user)

	return nil

}
