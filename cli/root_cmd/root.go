/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package rootcmd

import (
	"os"

	selfhost "github.com/atomic-blend/backend/cli/self-host"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ab",
	Short: "ab is the atomic blend cli tool",
	Long: `ab is a command-line tool for managing atomic blend applications.
	
It provides various commands and options to interact with the atomic blend platform.
Provides all the necessary tools and commands to self-host and manage atomic blend instances.
Also allows developeprs to manage their applications and services efficiently from the terminal.`,
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
	rootCmd.AddCommand(selfhost.NewCommand())
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
