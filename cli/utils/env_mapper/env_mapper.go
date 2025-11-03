package envmapper

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func MapFlagToEnv(cmd *cobra.Command, flagName, envName, configKey string) error {
	// Replace dashes with underscores and convert to uppercase
	err := viper.BindPFlag(configKey, cmd.Flags().Lookup(flagName))
	if err != nil {
		return err
	}
	err = viper.BindEnv(configKey, envName)
	if err != nil {
		return err
	}
	return nil
}
