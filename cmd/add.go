package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// UnstagedShortCode git status --short ouput XY code of not staged files
var UnstagedShortCode = []string{
	" A", " M", " D", "??",
}

// @return []string unstaged filepaths
func findUnstaged() {
	cmd := exec.Command(GitExec, "status", "--short")
	var out bytes.Buffer
	cmd.Stdout = &out
	must(cmd.Run())

	// scanner by default read by newlines
	scanner := bufio.NewScanner(&out)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		for _, code := range UnstagedShortCode {
			if strings.HasPrefix(line, code) {
				lines = append(lines, line)
			}
		}
	}

	// TODO: handle XY ORIG_PATH -> PATH
	// TODO: parse status from shortcode as well
	filepaths := []string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			filepaths = append(filepaths, fields[1])
		}
	}
	fmt.Println(filepaths)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "stage changes",
	Long:  "stage changes with prompt and select",
	Run: func(cmd *cobra.Command, args []string) {
		if mustBool(cmd.LocalFlags().GetBool("all")) {
			GitAddAll()
			return
		}
		// filepaths := findUnstaged()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("all", "a", false, "stage all changes")
}
