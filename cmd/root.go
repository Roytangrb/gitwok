package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Roytangrb/gitwok/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	isVerbose bool
)

var logger *util.Logger = util.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gitwok",
	Version: "v0.0.0",
	Short:   "Configurable CLI with conventional commits, changelog, git hooks all in one",
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
		logger.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	// global flags and configuration settings.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is gitwok.yaml)")
	rootCmd.PersistentFlags().BoolVar(&isVerbose, "verbose", false, "run commands with verbose output")

	// local flags and configuration settings.
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
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
			logger.Error(err)
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
			logger.Info("Using config file", viper.ConfigFileUsed())
		}
	} else {
		if fnfe, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error(fnfe.Error())
		} else if pe, ok := err.(viper.ConfigParseError); ok {
			logger.Error(pe.Error())
		} else {
			logger.Warn("Config file was found but another error was produced")
			logger.Error(err.Error())
		}
	}
}
