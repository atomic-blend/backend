// Package rootcmd contains the root command for the CLI.
//
// Copyright © 2025 NAME HERE <EMAIL ADDRESS>
package rootcmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/atomic-blend/backend/ab-cli/cmd/developper"
	selfhost "github.com/atomic-blend/backend/ab-cli/cmd/self-host"
	"github.com/atomic-blend/backend/ab-cli/config"
	envmapper "github.com/atomic-blend/backend/ab-cli/utils/viperutils"
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
	rootCmd.PersistentFlags().StringP("directory", "d", ".", "Directory where configurations and data are stored")
	envmapper.MapFlagToEnv(rootCmd, "directory", "ATOMIC_BLEND_DIRECTORY", "directory")
	rootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug logging")
	envmapper.MapFlagToEnv(rootCmd, "debug", "ATOMIC_BLEND_DEBUG", "debug")
	rootCmd.PersistentFlags().StringP("channel", "C", "", "Update channel to use (stable or rc)")
	envmapper.MapFlagToEnv(rootCmd, "channel", "ATOMIC_BLEND_CHANNEL", "channel")
	rootCmd.PersistentFlags().String("developper-mode", "", "Enable developper mode")
	envmapper.MapFlagToEnv(rootCmd, "developper-mode", "ATOMIC_BLEND_DEVELOPPER_MODE", "developper_mode")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /.ab-config.yaml)")

	rootCmd.AddCommand(selfhost.NewCommand())
	rootCmd.AddCommand(developper.NewCommand())

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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// 2. Bind Cobra flags to Viper early so flags (and their defaults) and env vars
	// are available when we need to determine paths (e.g. the --directory flag
	// is used to locate the config file). Bind both the command-local flags and
	// the root persistent flags so flags passed to subcommands are recognized.
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}
	// Also ensure root-level persistent flags are bound (covers some cobra versions
	// where cmd.PersistentFlags() may not include parent persistent flags).
	if root := cmd.Root(); root != nil {
		if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
			return err
		}
	}

	// 3. Configure Viper to read from the config file.
	// If --directory is set, use it; otherwise default to current directory.
	dir := viper.GetString("directory")
	if dir == "" {
		dir = "."
	}

	// Prefer explicit config file override via --config flag.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Prefer a hidden file named .ab-config.yaml if present in the directory.
		abConfigPath := filepath.Join(dir, ".ab-config.yaml")
		if _, err := os.Stat(abConfigPath); err == nil {
			viper.SetConfigFile(abConfigPath)
		} else {
			// Fall back to searching for a file named .ab-config.(yaml|yml|json) in the home directory.
			viper.AddConfigPath("$HOME")
			viper.SetConfigName(".ab-config")
			viper.SetConfigType("yaml")
		}
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

	// 5. Unmarshal the configuration into the CliConfig struct.
	err := viper.Unmarshal(&config.CliConfig)
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
