package cmd

import "testing"

// Testing if commitmsg.tmpl format the commit msg correctly
// trailing linebreak is expected
func TestCommitToString(t *testing.T) {
	case1 := CommitMsg{Type: "fix", Description: "hello"}

	// test are run inside package, use path relative to this file
	result := case1.ToString("../templates/commitmsg.tmpl")
	if result != "fix: hello\n" {
		t.Errorf(`commit msg format failed, expected: %s, got: %s`, `fix: hello\n`, result)
	}
}
