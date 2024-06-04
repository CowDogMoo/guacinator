package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cowdogmoo/guacinator/pkg/config"
	log "github.com/cowdogmoo/guacinator/pkg/logging"
	"github.com/l50/goutils/v2/sys"
	"github.com/mitchellh/go-homedir"
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
	guacConfigDir    string
	guacConfigFile   string
	cfg              config.Config

	debug bool

	rootCmd = &cobra.Command{
		Use:   "guacinator",
		Short: "Command line utility to interact programmatically with Apache Guacamole.",
	}

	home, _          = homedir.Dir()
	defaultConfigDir = filepath.Join(home, ".guacinator")
)

func init() {
	cobra.OnInitialize(initConfig)
	setupRootCmd(rootCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType(defaultConfigType)
	viper.AutomaticEnv()

	home, err := homedir.Dir()
	checkErr(err, "Failed to get home directory: %v")

	guacConfigDir = filepath.Join(home, defaultConfigDir)
	guacConfigFile = filepath.Join(guacConfigDir, fmt.Sprintf("%s.%s", defaultConfigName, defaultConfigType))

	// Check if the config file exists, if not create the default config file
	if _, err := os.Stat(guacConfigDir); os.IsNotExist(err) {
		fmt.Printf("Config file not found, creating default config file at %s", guacConfigFile)
		createConfig(guacConfigFile)
	}

	viper.SetConfigFile(guacConfigFile)

	if err := viper.ReadInConfig(); err != nil {
		checkErr(err, "Can't read config: %v")
	}

	if err := viper.Unmarshal(&guacCfg); err != nil {
		checkErr(err, "Failed to unmarshal config: %v")
	}

	err = log.Initialize(guacConfigDir, cfg.Log.Level, cfg.Log.LogPath)
	checkErr(err, "Failed to initialize the logger: %v")

	// Check for required dependencies after initializing the logger
	checkErr(depCheck(), "Dependency check failed")
}

func createConfig(cfgPath string) {
	cfgDir := filepath.Dir(cfgPath)

	// Ensure the configuration directory exists
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		fmt.Printf("Creating config directory %s", cfgDir)
		checkErr(os.MkdirAll(cfgDir, os.ModePerm), "failed to create config directory %s: %v")
	}

	// Write the default config file if it does not exist
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		configFileData, err := configContentsFs.ReadFile(filepath.Join("config", "config.yaml"))
		checkErr(err, "failed to read embedded config: %v")
		checkErr(os.WriteFile(cfgPath, configFileData, 0644), "failed to write config to %s: %v")
		fmt.Printf("Default config file created at %s", cfgPath)
	} else {
		fmt.Printf("Config file already exists at %s", cfgPath)
	}
}

func setupRootCmd(cmd *cobra.Command) {
	pf := cmd.PersistentFlags()
	pf.StringVar(&guacConfigFile, "config", "", "config file (default is $HOME/.guacinator/guacinator-config.yaml)")
	if err := viper.BindPFlag("config", pf.Lookup("config")); err != nil {
		log.Error("Failed to bind the config flag: %v", err)
	}

	pf.BoolVarP(
		&debug, "debug", "d", false, "Show debug messages.")
	if err := viper.BindPFlag("debug", pf.Lookup("debug")); err != nil {
		checkErr(err, "Failed to bind the debug flag: %v")
	}

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkErr(err error, format string) {
	if err != nil {
		log.Error(format, err)
		os.Exit(1)
	}
}

func depCheck() error {
	if !sys.CmdExists("kubectl") {
		errMsg := "required program kubectl is not installed in $PATH, exiting"
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	log.Debug("All dependencies are satisfied.")

	return nil
}

// Execute runs the root cobra command. It checks for errors and exits
// the program if any are encountered.
func Execute() {
	checkErr(rootCmd.Execute(), "Command execution failed")
}
