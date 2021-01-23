package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GitExec git executable name
const GitExec = "git"

// Git with methods to exec git commands
type Git struct {
	verbose bool
	dryRun  bool
}

func hasDryRunFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-n" || arg == "--dry-run" {
			return true
		}
	}
	return false
}

func prependArg(arg string, args []string) []string {
	return append([]string{arg}, args...)
}

// Add exec `git add <args>`
func (git *Git) Add(args ...string) {
	if !hasDryRunFlag(args) && git.dryRun {
		args = prependArg("--dry-run", args)
	}
	cmd := exec.Command(GitExec, prependArg("add", args)...)
	must(cmd.Run())
}

// Rm exec `git rm` to stage changes of a deleted file
func (git *Git) Rm(args ...string) {
	if !hasDryRunFlag(args) && git.dryRun {
		args = prependArg("--dry-run", args)
	}
	cmd := exec.Command(GitExec, prependArg("rm", args)...)
	must(cmd.Run())
}

// GitCommit exec `git commit -m`
func GitCommit(msg string) {
	cmd := exec.Command(GitExec, "commit", "-m", msg)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	// print output, delay exit on error
	logger.Info("git commit -m output:")
	fmt.Println(out.String())

	must(err)
}
