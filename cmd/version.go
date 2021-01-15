package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version",
	Long:  "print and its dependencies' versions",
	Run: func(cmd *cobra.Command, args []string) {
		printlnVer(rootCmd.Name(), rootCmd.Version)
		fmt.Println()
		printlnVer("Conventional-commits-spec", "v1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printlnVer(name string, semver string) {
	fmt.Printf("%s %s\n", name, semver)
}
