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

// git status --short (more detail)
// In the short-format, the status of each path is shown as one of these forms
//    XY PATH
//    XY ORIG_PATH -> PATH
// X shows the status of the index, and Y shows the status of the work tree.

// TODO: handle merge conflict cases for index status

const (
	// CodeAddedNotStaged added not staged
	CodeAddedNotStaged = " A"
	// CodeModifiedNotStaged modified not staged
	CodeModifiedNotStaged = " M"
	// CodeDeletedNotStaged deleted not staged
	CodeDeletedNotStaged = " D"
	// CodeRenamedNotStaged renamed in work tree not staged
	// not common, `git rename` will automatically stage the changes
	CodeRenamedNotStaged = " R"
	// CodeCopiedNotStaged copied in work tree not staged
	// not common, TODO: verify usage
	CodeCopiedNotStaged = " C"
	// CodeUntracked untracked
	CodeUntracked = "??"
	// PathSepArrow path separator if file is moved or copied
	PathSepArrow = "->"
)

func translateNotStaged(code string) string {
	switch code {
	case CodeAddedNotStaged:
		return "added"
	case CodeModifiedNotStaged:
		return "modified"
	case CodeDeletedNotStaged:
		return "deleted"
	case CodeRenamedNotStaged:
		return "renamed"
	case CodeCopiedNotStaged:
		return "copied"
	case CodeUntracked:
		return "untracked"
	default:
		return "unknown"
	}
}

// UnstagedShortCodes git status --short output XY code of not staged files
var UnstagedShortCodes = []string{
	CodeAddedNotStaged,
	CodeModifiedNotStaged,
	CodeDeletedNotStaged,
	CodeRenamedNotStaged,
	CodeCopiedNotStaged,
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
		for _, code := range UnstagedShortCodes {
			if strings.HasPrefix(line, code) {
				lines = append(lines, line)
			}
		}
	}

	codes := []string{}
	filepaths := []string{}
	for _, line := range lines {
		code := line[:2]
		rest := strings.TrimSpace(line[2:])
		if strings.Contains(rest, PathSepArrow) {
			// TODO: handle file renamed or copied ORIG_PATH -> PATH
		} else {
			codes = append(codes, code)
			// filepath whitespace no encoded, should be handled in later steps
			filepaths = append(filepaths, rest)
		}
	}

	return codes, filepaths
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "stage changes",
	Long:  "stage changes with prompt and select",
	Run: func(cmd *cobra.Command, args []string) {
		var git = &Git{
			verbose: false,
			dryRun:  mustBool(cmd.Flags().GetBool("dry-run")),
		}

		if mustBool(cmd.LocalFlags().GetBool("all")) {
			git.Add(".")
			return
		}

		codes, filepaths := findUnstaged()
		if len(filepaths) > 0 {
			codeDict := make(map[string]string)
			fpDict := make(map[string]string)

			labels := []string{}
			for i, fp := range filepaths {
				code := codes[i]
				label := translateNotStaged(code) + ": " + fp
				labels = append(labels, label)
				codeDict[label] = code
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
				code := codeDict[label]
				fp := fpDict[label]

				if code == CodeDeletedNotStaged {
					git.Rm(fp)
				} else {
					if _, err := os.Stat(fp); err == nil {
						git.Add(fp)
					} else {
						logger.Warn(err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("all", "a", false, "stage all changes")
}
