package cmd

import "testing"

// TestStr string test helper struct
// put err field the last as optional
type TestStr struct {
	got      string
	expected string
	msg      string
}

func TestTranslateNotStaged(t *testing.T) {
	var tests = []TestStr{
		{translateNotStaged(CodeAddedNotStaged), "added", ""},
		{translateNotStaged(CodeModifiedNotStaged), "modified", ""},
		{translateNotStaged(CodeDeletedNotStaged), "deleted", ""},
		{translateNotStaged(CodeRenamedNotStaged), "renamed", ""},
		{translateNotStaged(CodeCopiedNotStaged), "copied", ""},
		{translateNotStaged(CodeUntracked), "untracked", ""},
		{translateNotStaged("M "), "unknown", ""},
		{translateNotStaged("D "), "unknown", ""},
	}

	for _, test := range tests {
		if test.got != test.expected {
			t.Errorf("translateNotStaged failed, expected: %s, got: %s", test.expected, test.got)
		}
	}
}
