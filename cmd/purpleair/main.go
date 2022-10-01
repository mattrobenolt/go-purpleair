package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.withmatt.com/purpleair/cmd/purpleair/internal/addsensor"
	"go.withmatt.com/purpleair/cmd/purpleair/internal/keys"
	"go.withmatt.com/purpleair/cmd/purpleair/internal/sensors"
)

var (
	rootCmd = &cobra.Command{
		Use: "purpleair",
	}

	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.PersistentFlags()
	flags.StringVar(&configFile, "config", "", "config file (default is $HOME/.purpleair.yaml)")
	flags.String("read-key", "", "API read key")
	flags.String("write-key", "", "API write key")

	viper.BindPFlag("read-key", flags.Lookup("read-key"))
	viper.BindPFlag("write-key", flags.Lookup("write-key"))

	addsensor.Init(rootCmd)
	keys.Init(rootCmd)
	sensors.Init(rootCmd)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".purpleair")
	}

	viper.SetEnvPrefix("PURPLEAIR")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cobra.CheckErr(err)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
