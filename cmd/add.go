package cmd

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

const (
	// CodeAddedNotStaged added not staged
	CodeAddedNotStaged = " A"
	// CodeModifiedNotStaged modified not staged
	CodeModifiedNotStaged = " M"
	// CodeDeletedNotStaged deleted not staged
	CodeDeletedNotStaged = " D"
	// CodeUntracked untracked
	CodeUntracked = "??"
	// CodeRenamedNotStaged renamed in work tree not staged
	CodeRenamedNotStaged = " R" // TODO: verify meaning
	// CodeCopiedNotStaged copied in work tree not staged
	CodeCopiedNotStaged = " C" // TODO: verify meaning
)

func translateShortCode(code string) string {
	switch code {
	case CodeAddedNotStaged:
		return "added"
	case CodeModifiedNotStaged:
		return "modified"
	case CodeDeletedNotStaged:
		return "deleted"
	case CodeUntracked:
		return "untracked"
	case CodeRenamedNotStaged:
		return "renamed"
	case CodeCopiedNotStaged:
		return "copied"
	default:
		return "unknown"
	}
}

// UnstagedShortCode git status --short output XY code of not staged files
var UnstagedShortCode = []string{
	CodeAddedNotStaged,
	CodeModifiedNotStaged,
	CodeDeletedNotStaged,
	CodeUntracked,
}

// @return []string codes unstaged file short code status
// @return []string filepaths unstaged filepaths
func findUnstaged() ([]string, []string) {
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
	codes := []string{}
	filepaths := []string{}
	for _, line := range lines {
		code := line[:2]
		rest := strings.TrimSpace(line[2:])
		if !ContainsWhiteSpace(rest) {
			codes = append(codes, code)
			filepaths = append(filepaths, rest)
		} else {
			// TODO: handle filepath with whitespace

		}
	}

	return codes, filepaths
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

		codes, filepaths := findUnstaged()
		if len(filepaths) > 0 {
			fpDict := make(map[string]string)
			labels := []string{}
			for i, fp := range filepaths {
				code := codes[i]
				label := translateShortCode(code) + ": " + fp
				labels = append(labels, label)
				fpDict[label] = fp
			}

			selectedLabels := []string{}
			prompt := &survey.MultiSelect{
				Message: "Stage changes to commit:",
				Options: labels,
			}
			// exit on prompt interrupted
			must(survey.AskOne(prompt, &selectedLabels))

			for _, label := range selectedLabels {
				fp := fpDict[label]
				// verify filepath
				if _, err := os.Stat(fp); err == nil {
					GitAdd(fp)
				} else {
					logger.Warn(err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("all", "a", false, "stage all changes")
}
