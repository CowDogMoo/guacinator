package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/l50/goutils/v2/logging"
	"github.com/l50/goutils/v2/sys"

	"github.com/mitchellh/go-homedir"

	"golang.org/x/exp/slog"

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
	logger  logging.Logger
	logFile *os.File

	home, _          = homedir.Dir()
	defaultConfigDir = filepath.Join(home, ".guacinator")
)

// Execute runs the root cobra command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
	defer logFile.Close()
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err = homedir.Dir()
	if err != nil {
		logger.Println("failed to get the home directory:", err)
		cobra.CheckErr(err)
	}

	pf := rootCmd.PersistentFlags()
	pf.StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.guacinator/guacinator-config.yaml)")

	pf.BoolVarP(
		&debug, "debug", "", false, "Show debug messages.")

	if err := viper.BindPFlag("debug", pf.Lookup("debug")); err != nil {
		fmt.Printf("Failed to get the home directory: %v\n", err)
		cobra.CheckErr(err)
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configLogging() error {
	// Set log levels
	configLogLevel := viper.GetString("log.level")
	var level slog.Level
	switch configLogLevel {
	case "debug":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	// Create log file handlers
	logFilePath := filepath.Join(defaultConfigDir, "guacinator.log")
	logger, err = logging.ConfigureLogger(level, logFilePath)
	if err != nil {
		return fmt.Errorf("failed to set up logging: %w", err)
	}

	// Logging message
	logger.Println("Initialization complete! Logging setup successfully.")

	return nil
}

func getConfigFile() ([]byte, error) {
	configFileData, err := configContentsFs.ReadFile(
		filepath.Join("config", defaultConfigName))
	if err != nil {
		logger.Error("error reading config/ contents: %v", err)
		return configFileData, err
	}

	return configFileData, nil
}

func createConfigFile(cfgPath string) error {
	if err := os.MkdirAll(filepath.Dir(cfgPath), os.ModePerm); err != nil {
		logger.Error("cannot create dir %s: %s", cfgPath, err)
		return err
	}

	configFileData, err := getConfigFile()
	if err != nil {
		logger.Error("failed to get config file data: %v", err)
		return err
	}

	if err := os.WriteFile(cfgPath, configFileData, os.ModePerm); err != nil {
		logger.Error("failed to write to config file: %v\n", err)
		return err
	}

	if err := os.WriteFile(cfgPath, configFileData, os.ModePerm); err != nil {
		logger.Error("failed to write to config file: %v", err)
		return err
	}

	cmd := "kubectl"
	if !sys.CmdExists(cmd) {
		err := fmt.Errorf("required program %s is not installed in $PATH, exiting", cmd)
		logger.Error(err.Error())
		return err
	}

	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if err := configLogging(); err != nil {
		fmt.Printf("failed to set up logging: %v\n", err)
		cobra.CheckErr(err)
	}

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
		logger.Println(color.BlueString(
			"No config file found - creating " +
				filepath.Join(defaultConfigDir,
					defaultConfigName) +
				" with default values"))

		if err := createConfigFile(
			filepath.Join(defaultConfigDir, defaultConfigName)); err != nil {
			logger.Error("failed to create the config file: %v", err)
			cobra.CheckErr(err)
		}

		if err := viper.ReadInConfig(); err != nil {
			logger.Error("failed to read contents of config file: %v", err)
			cobra.CheckErr(err)
		} else {
			logger.Printf("Using config file: %s", viper.ConfigFileUsed())
		}
	} else {
		logger.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}
