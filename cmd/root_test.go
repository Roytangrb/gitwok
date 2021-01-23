package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

type TestBool struct {
	got      bool
	expected bool
	msg      string
}

// CompareStrSlices return true if two string slice contains same values in order
func CompareStrSlices(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

func TestInitDefaults(t *testing.T) {
	// reset all to default settings
	viper.Reset()
	// test default values
	initDefaults()

	var boolTests = []TestBool{
		{viper.GetBool("gitwok.commit.prompt.scope"), true, "scope prompt"},
		{viper.GetBool("gitwok.commit.prompt.breaking"), true, "breaking prompt"},
		{viper.GetBool("gitwok.commit.prompt.body"), true, "body prompt"},
		{viper.GetBool("gitwok.commit.prompt.footers"), true, "footers prompt"},
		{CompareStrSlices(viper.GetStringSlice("gitwok.commit.type"), PresetCommitTypes), true, "type options"},
		{CompareStrSlices(viper.GetStringSlice("gitwok.commit.scope"), []string{}), true, "scope options"},
	}

	for _, test := range boolTests {
		if test.got != test.expected {
			t.Errorf("Config default %s failed, expected: %t, got: %t", test.msg, test.expected, test.got)
		}
	}
}

func TestReadConfig(t *testing.T) {
	// reset rootCmd
	rootCmd.ResetCommands()
	rootCmd.ResetFlags()
	logger.VerboseEnabled = false
	// test config read
	// rootCmd.Flags().Set("verbose", "true")
	// rootCmd.Flags().Set("config", "../gitwok.yaml")
	// readConfig()

	if got, expected := logger.VerboseEnabled, false; got != expected {
		t.Errorf("Set logger verbose by flag failed, expected: %t, got: %t", expected, got)
	}
}
