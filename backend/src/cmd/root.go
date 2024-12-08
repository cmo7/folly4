package cmd

import (
	"fmt"
	"os"

	"github.com/cmo7/folly4/src/cmd/config"
	"github.com/cmo7/folly4/src/cmd/serve"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "folly",
	Short: "Una aplicación fullstack",
	Long:  `Una aplicación fullstack que incluye un backend en Go y un frontend en React.`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./conf)")
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(serve.ServeCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}

	//* Set defaults
	viper.SetDefault("app.name", "Folly")
	viper.SetDefault("app.version", "0.0.1")
}
