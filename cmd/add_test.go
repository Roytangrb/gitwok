package cmd

import (
	"bytes"
	"testing"
)

func TestTranslateNotStaged(t *testing.T) {
	var tests = []TestStr{
		{CodeAddedNotStaged, "added", ""},
		{CodeModifiedNotStaged, "modified", ""},
		{CodeDeletedNotStaged, "deleted", ""},
		{CodeRenamedNotStaged, "renamed", ""},
		{CodeCopiedNotStaged, "copied", ""},
		{CodeUntracked, "untracked", ""},
		{"M ", "unknown", ""},
		{"D ", "unknown", ""},
	}

	for _, test := range tests {
		if translateNotStaged(test.got) != test.expected {
			t.Errorf("translateNotStaged failed, expected: %s, got: %s", test.expected, test.got)
		}
	}
}

func TestFindUnstaged(t *testing.T) {
	emptyStatus := ``
	spaceStatus := ` M cmd/file with space.go`
	modifiedStatus := ` M cmd/add.go
 M cmd/add_test.go`

	var tests = []struct {
		in        string
		codes     []string
		filepaths []string
	}{
		{emptyStatus, []string{}, []string{}},
		{spaceStatus, []string{" M"}, []string{"cmd/file with space.go"}},
		{modifiedStatus, []string{" M", " M"}, []string{"cmd/add.go", "cmd/add_test.go"}},
	}

	for _, test := range tests {
		var out bytes.Buffer
		out.Write([]byte(test.in))
		codes, fps := findUnstaged(&out)

		if !CompareStrSlices(codes, test.codes) {
			t.Errorf("findUnstaged get codes failed, expected: %v, got: %v", test.codes, codes)
		}
		if !CompareStrSlices(fps, test.filepaths) {
			t.Errorf("findUnstaged get filepaths failed, expected: %v, got: %v", test.filepaths, fps)
		}
	}
}
