/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package rootcmd

import (
	"errors"
	"os"
	"strings"

	"github.com/atomic-blend/backend/cli/config"
	selfhost "github.com/atomic-blend/backend/cli/self-host"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ab",
	Short: "ab is the atomic blend cli tool",
	Long: `ab is a command-line tool for managing atomic blend applications.
	
It provides various commands and options to interact with the atomic blend platform.
Provides all the necessary tools and commands to self-host and manage atomic blend instances.
Also allows developeprs to manage their applications and services efficiently from the terminal.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /.ab-config.yaml)")
	rootCmd.AddCommand(selfhost.NewCommand())
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

func initializeConfig(cmd *cobra.Command) error {
	// 1. Set up Viper to use environment variables.
	viper.SetEnvPrefix("ATOMIC_BLEND")
	// Allow for nested keys in environment variables (e.g. `ATOMIC_BLEND_DATABASE_HOST`)
	// Replace dots and dashes with underscores so keys like `selfhost.directory`
	// map to env vars like `ATOMIC_BLEND_SELFHOST_DIRECTORY`.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "*"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for a config file with the name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(".ab-config.yaml")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 3. Read the configuration file.
	// If a config file is found, read it in. We use a robust error check
	// to ignore "file not found" errors, but panic on any other error.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// 4. Bind Cobra flags to Viper.
	// This is the magic that makes the flag values available through Viper.
	// It binds the full flag set of the command passed in.
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	// 5. Unmarshal the configuration into the CliConfig struct.
	err = viper.Unmarshal(&config.CliConfig)
	if err != nil {
		return err
	}

	debug := os.Getenv("LOG_LEVEL")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug == "debug" || config.CliConfig.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Interface("config", config.CliConfig).Msg("Configuration initialized successfully")

	return nil
}
