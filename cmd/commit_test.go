package cmd

import "testing"

// Testing if commitmsg.tmpl format the commit msg correctly
// trailing linebreak is expected
func TestCommitToString(t *testing.T) {
	case1 := CommitMsg{Type: "fix", Description: "hello"}

	result := case1.String()
	if result != "fix: hello\n" {
		t.Errorf(`commit msg format failed, expected: %s, got: %s`, `fix: hello\n`, result)
	}
}
