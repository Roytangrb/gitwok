package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print GitWok version",
	Long:  "Print GitWok and its dependencies' versions",
	Run: func(cmd *cobra.Command, args []string) {
		printVer()
		fmt.Println()
		printVerDeps()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVer() {
	fmt.Println("GitWok v0.0.0")
}

func printVerDeps() {
	fmt.Println("Conventional-commits-spec: v1.0.0")
}
