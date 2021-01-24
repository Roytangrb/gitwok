package cmd

import "testing"

func TestHasDryRunFlag(t *testing.T) {
	var tests = []TestBool{
		{hasDryRunFlag([]string{}), false, "empty flags"},
		{hasDryRunFlag([]string{"--dry-run"}), true, "--dry-run flag"},
		{hasDryRunFlag([]string{"-n"}), false, "-n flag"},
	}

	for _, test := range tests {
		if test.got != test.expected {
			t.Errorf("TestHasDryRunFlag with %s failed, expected: %t, got: %t", test.msg, test.expected, test.got)
		}
	}
}

func TestPrependArg(t *testing.T) {
	if got, expected := prependArg("--dry-run", []string{"add"}), []string{"--dry-run", "add"}; !CompareStrSlices(got, expected) {
		t.Errorf("TestPrependArg failed, expected: %v, got: %v", expected, got)
	}
}
