package cmd

import "testing"

type TestStrInOut struct {
	in  string
	out string
}

func TestTranslateNotStaged(t *testing.T) {
	var tests = []TestStrInOut{
		{CodeAddedNotStaged, "added"},
		{CodeModifiedNotStaged, "modified"},
		{CodeDeletedNotStaged, "deleted"},
		{CodeRenamedNotStaged, "renamed"},
		{CodeCopiedNotStaged, "copied"},
		{CodeUntracked, "untracked"},
		{"M ", "unknown"},
		{"D ", "unknown"},
	}

	for _, test := range tests {
		if got, expected := translateNotStaged(test.in), test.out; got != expected {
			t.Errorf("translateNotStaged failed, expected: %s, got: %s", expected, got)
		}
	}
}
