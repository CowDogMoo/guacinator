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
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	goutils "github.com/l50/goutils"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "guacinator-config.yaml"
	defaultConfigType = "yaml"
)

var (
	//go:embed config/*
	configContentsFs embed.FS

	cfgFile string
	debug   bool
	err     error

	rootCmd = &cobra.Command{
		Use:   "guacinator",
		Short: "Command line utility to interact programmatically with Apache Guacamole.",
	}

	home, _          = homedir.Dir()
	defaultConfigDir = filepath.Join(home, ".guacinator")
)

// Execute runs the root cobra command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err = homedir.Dir()
	if err != nil {
		os.Exit(1)
	}

	pf := rootCmd.PersistentFlags()
	pf.StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.guacinator/guacinator-config.yaml)")

	pf.BoolVarP(
		&debug, "debug", "", false, "Show debug messages.")

	if err := viper.BindPFlag("debug", pf.Lookup("debug")); err != nil {
		log.WithError(err).Error("failed to bind to debug in the yaml config")
		os.Exit(1)
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configLogging() error {
	logger, err := goutils.CreateLogFile()
	if err != nil {
		log.WithError(err).Error("error creating the log file")
	}

	// Set log level
	configLogLevel := viper.GetString("log.level")
	if logLevel, err := log.ParseLevel(configLogLevel); err != nil {
		log.WithFields(log.Fields{"level": logLevel,
			"fallback": "info"}).Warn("Invalid log level")
	} else {
		if debug {
			log.Debug("Debug logs enabled")
			logLevel = log.DebugLevel
		}
		log.SetLevel(logLevel)
	}

	// Set log format
	switch viper.GetString("log.format") {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		})
	}

	// Output to both stdout and the log file
	mw := io.MultiWriter(os.Stdout, logger.FilePtr)
	log.SetOutput(mw)

	return nil
}

func getConfigFile() ([]byte, error) {
	configFileData, err := configContentsFs.ReadFile(
		filepath.Join("config", defaultConfigName))
	if err != nil {
		log.WithError(err).Errorf("error reading config/ contents: %v", err)
		return configFileData, err
	}

	return configFileData, nil
}

func createConfigFile(cfgPath string) error {
	if err := os.MkdirAll(filepath.Dir(cfgPath), os.ModePerm); err != nil {
		log.WithError(err).Errorf("cannot create dir %s: %s", cfgPath, err)
		return err
	}

	configFileData, err := getConfigFile()
	if err != nil {
		log.WithError(err).Errorf("cannot get lines of config file: %v", err)
		return err
	}

	if err := os.WriteFile(cfgPath, configFileData, os.ModePerm); err != nil {
		log.WithError(err)
		return err
	}

	cmd := "kubectl"
	if !goutils.CmdExists(cmd) {
		err := fmt.Errorf("required program %s is not installed in $PATH, exiting", cmd)
		log.WithError(err)
		return err
	}

	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config yaml file in the config directory
		viper.AddConfigPath(defaultConfigDir)
		viper.SetConfigType(defaultConfigType)
		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Info(color.BlueString(
			"No config file found - creating " +
				filepath.Join(defaultConfigDir,
					defaultConfigName) +
				" with default values"))

		if err := createConfigFile(
			filepath.Join(defaultConfigDir, defaultConfigName)); err != nil {
			log.WithError(err).Error("failed to create the config file")
			os.Exit(1)
		}

		if err := viper.ReadInConfig(); err != nil {
			log.WithError(err).Error("error reading config file")
			os.Exit(1)
		} else {
			log.Debug("Using config file: ", viper.ConfigFileUsed())
		}
	} else {
		log.Debug("Using config file: ", viper.ConfigFileUsed())
	}

	if err := configLogging(); err != nil {
		log.WithError(err).Error("failed to set up logging")
	}
}
