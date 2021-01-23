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

// hasDryRunFlag check if "--dry-run" is passed in already
// caveats: "-n" is not checked because it could mean sth else
// i.e. in git commit "-n" means "--no-verify"
// "--dry-run" should be passed in instead of "-n"
func hasDryRunFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--dry-run" {
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

// Commit exec `git commit <args>`
func (git *Git) Commit(args ...string) {
	if !hasDryRunFlag(args) && git.dryRun {
		args = prependArg("--dry-run", args)
	}
	cmd := exec.Command(GitExec, prependArg("commit", args)...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	// print output, delay exit on error
	logger.Verbose(fmt.Sprintln("git commit output:") + out.String())

	must(err)
}
