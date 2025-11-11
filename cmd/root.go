package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hscli",
	Short: "A CLI tool for managing HubSpot contacts",
	Long: `hscli is a command-line tool for managing HubSpot contacts.
It provides CRUD operations for contacts including listing, creating,
updating, deleting, and querying contacts.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hscli.yaml)")
	rootCmd.PersistentFlags().String("api-key", "", "HubSpot API key (or set HUBSPOT_API_KEY env var)")
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hscli")
	}

	viper.SetEnvPrefix("HUBSPOT")
	viper.AutomaticEnv() // read in environment variables that match

	// Bind the api-key to the HUBSPOT_API_KEY environment variable
	viper.BindEnv("api-key", "HUBSPOT_API_KEY")

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
