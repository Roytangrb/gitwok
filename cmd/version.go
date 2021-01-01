package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var isPrintVerbose bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print GitWok version",
	Long:  "Print GitWok and its dependencies' versions",
	Run: func(cmd *cobra.Command, args []string) {
		printVer()

		if isPrintVerbose {
			printVerDeps()
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&isPrintVerbose, "verbose", "v", false, "Print dependencies versions")
}

func printVer() {
	fmt.Println("GitWok v0.0.0")
}

func printVerDeps() {
	fmt.Println("Dependencies: Not Implemented")
}
