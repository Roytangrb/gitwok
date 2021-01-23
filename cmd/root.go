package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Roytangrb/gitwok/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// VersionTmpl version template for --version output
const VersionTmpl = `
{{- .Name}} {{.Version}}
`

var logger *util.Logger = util.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gitwok",
	Version: "v0.2.0",
	Short:   "Configurable CLI with conventional commits, changelog, git hooks all in one",
	Run:     func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	must(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initDefaults, readConfig)

	rootCmd.SetVersionTemplate(VersionTmpl)

	rootCmd.PersistentFlags().String("config", "", "config file (default is ./gitwok.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "dry run all git exec actions")
}

func initDefaults() {
	viper.SetDefault("gitwok.commit.prompt.scope", true)
	viper.SetDefault("gitwok.commit.prompt.breaking", true)
	viper.SetDefault("gitwok.commit.prompt.body", true)
	viper.SetDefault("gitwok.commit.prompt.footers", true)
	viper.SetDefault("gitwok.commit.type", PresetCommitTypes)
	viper.SetDefault("gitwok.commit.scope", []string{})
}

func readConfig() {
	logger.VerboseEnabled = mustBool(rootCmd.Flags().GetBool("verbose"))

	if fp := mustStr(rootCmd.Flags().GetString("config")); fp != "" {
		// Use config file from the flag.
		viper.SetConfigFile(fp)
	} else {
		// Find home directory.
		home := mustStr(homedir.Dir())
		// Search config in cwd or home directory
		viper.SetConfigName("gitwok")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	// TODO: read in environment variables that match
	// viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logger.Verbose("Using config file", viper.ConfigFileUsed())
	} else {
		if fnfe, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn(fnfe)
		} else if pe, ok := err.(viper.ConfigParseError); ok {
			logger.Warn(pe)
		} else {
			logger.Warn(err)
		}
	}
}

func must(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func mustStr(val string, err error) string {
	must(err)
	return val
}

func mustStrSlice(val []string, err error) []string {
	must(err)
	return val
}

func mustBool(val bool, err error) bool {
	must(err)
	return val
}
