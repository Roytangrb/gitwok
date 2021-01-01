package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	isVerbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitwok",
	Short: "Configurable CLI with conventional commits, changelog, git hooks all in one",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global flags and configuration settings.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is gitwok.yaml)")
	rootCmd.PersistentFlags().BoolVar(&isVerbose, "verbose", false, "run commands with verbose output")

	// local flags and configuration settings.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in cwd or home directory
		viper.SetConfigName("gitwok")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if isVerbose {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if fnfe, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("[Error] ", fnfe.Error())
		} else if pe, ok := err.(viper.ConfigParseError); ok {
			fmt.Println("[Error] ", pe.Error())
		} else {
			fmt.Println("Config file was found but another error was produced")
			fmt.Println(err.Error())
		}
	}
}
