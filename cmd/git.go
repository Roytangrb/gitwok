package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GitExec git executable name
const GitExec = "git"

// GitAddAll exec `git add .`
func GitAddAll() {
	cmd := exec.Command(GitExec, "add", ".")
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
