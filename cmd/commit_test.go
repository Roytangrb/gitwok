package cmd

import (
	"fmt"
	"strings"
	"testing"
)

var NL = fmt.Sprintln()

func TestContainsNewline(t *testing.T) {
	if !ContainsNewline(fmt.Sprintln()) {
		t.Error("ContainsNewline check failed")
	}
}

// TestMakeCommitMsg test trim string input
// - whitespaces including `\r\n` on Windows should be trimmed
func TestMakeCommitMsg(t *testing.T) {
	msg := makeCommitMsg(" fix ", " lib ", false, " desc\r\n", " bodyln1\n bodyln2\n", []string{" Acked-by: RT \n", "Review-by: RT2\r\n", " #1", "Reviewed: "})

	var tests = []TestStrInOut{
		{msg.Type, "fix"},
		{msg.Scope, "lib"},
		{msg.Description, "desc"},
		{msg.Body, "bodyln1\n bodyln2"},
		{msg.Footers[0], "Acked-by: RT"},
		{msg.Footers[1], "Review-by: RT2"},
		{msg.Footers[2], " #1"},
		{msg.Footers[3], "Reviewed: "},
	}

	for _, test := range tests {
		if got, expected := test.in, test.out; got != expected {
			t.Errorf("makeCommitMsg failed, expected: %q, got: %q", expected, got)
		}
	}
}

// TestCommitMsgHeader
// @spec conventional commits v1.0.0
// 1. REQUIRED `type` of a noun, OPTIONAL `scope`, OPTIONAL `!``, and REQUIRED terminal colon and space.
// 4. `scope` MUST consist of a noun describing a section of the codebase surrounded by parenthesis
// 5. `description` MUST immediately follow the colon and space after the type/scope prefix.
func TestCommitMsgHeader(t *testing.T) {
	// test CommitMsg.Validate
	// - required `type` and `description`
	// - no linebreaks within `type`, `scope`, and `description`
	emptyTypeMsg := makeCommitMsg("", "", false, "", "", []string{})
	invalidTypeMsg := makeCommitMsg("a type", "", false, "desc", "", []string{})
	invalidScopeMsg := makeCommitMsg("fix", "with"+NL+"linebreak", false, "", "", []string{})
	emptyDescMsg := makeCommitMsg("fix", "", false, "", "", []string{})
	invalidDescMsg := makeCommitMsg("fix", "", false, "desc with line"+NL+"line2", "", []string{})

	if ok, msg := emptyTypeMsg.Validate(); ok || msg != RequiredType {
		t.Errorf("Required commit type check failed, expected: %q, got: %q", RequiredType, msg)
	}

	if ok, msg := invalidTypeMsg.Validate(); ok || msg != InvalidType {
		t.Errorf("Invalid commit type check failed, expected: %q, got: %q", InvalidType, msg)
	}

	if ok, msg := invalidScopeMsg.Validate(); ok || msg != InvalidScope {
		t.Errorf("No linebreak in scope check failed, expected: %q, got: %q", InvalidScope, msg)
	}

	if ok, msg := emptyDescMsg.Validate(); ok || msg != RequiredDesc {
		t.Errorf("Required commit description check failed, expected: %q, got: %q", RequiredDesc, msg)
	}

	if ok, msg := invalidDescMsg.Validate(); ok || msg != InvalidDesc {
		t.Errorf("Invalid commit description check failed, expected: %q, got: %q", InvalidDesc, msg)
	}

	// test CommitMsg.ToString
	var tests = []TestStrInOut{
		{makeCommitMsg(" docs ", "", false, "fix typo", "", []string{}).ToString(), "docs: fix typo" + NL},                       // test type trim
		{makeCommitMsg("docs", " READ ME.md ", false, "fix typo", "", []string{}).ToString(), "docs(READ ME.md): fix typo" + NL}, // test scope trim
		{makeCommitMsg("docs", "", true, "fix typo", "", []string{}).ToString(), "docs!: fix typo" + NL},
		{makeCommitMsg("fix", "lib", true, "fix bug", "", []string{}).ToString(), "fix(lib)!: fix bug" + NL},
	}

	for _, test := range tests {
		if got, expected := test.in, test.out; got != expected {
			t.Errorf("CommitMsg.ToString for header failed, expected: %q, got: %q", expected, got)
		}
	}
}

// TestCommitMsgBody
// @spec conventional commits v1.0.0
// 6. body MUST begin one blank line after the description
func TestCommitMsgBody(t *testing.T) {
	body := fmt.Sprintln("msg body") + NL + "body line2"
	msg1 := makeCommitMsg("docs", "", false, "fix typo", body, []string{})

	if got, expected := msg1.ToString(), fmt.Sprintln("docs: fix typo")+NL+fmt.Sprintln(body); got != expected {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}

func TestFootersWriteAnswer(t *testing.T) {
	var cmtFooters CommitFooters
	if err := cmtFooters.WriteAnswer("footers", 1); err == nil {
		t.Errorf("Write answer to custom footers struct failed, reason: asserting string input should return error")
	}

	if err := cmtFooters.WriteAnswer("footers", ""); err != nil || len(cmtFooters.Footers) != 0 {
		t.Errorf("Write answer to custom footers struct failed, reason: empty input should return []string")
	}

	if err := cmtFooters.WriteAnswer("footers", "Acked-By: RT"); err != nil {
		t.Error("Write answer to custom footers struct failed, got error", err)
	}
}

func TestMatchFooters(t *testing.T) {
	if footers := MatchFooters(""); len(footers) != 0 {
		t.Errorf("MatchFooters failed, empty input expected matched %d, matched %d, %v", 0, len(footers), footers)
	}

	var testStr = `
Acked-By: RT fix readme
with second line

Reviewed-By: RT
fix #1
		
  BREAKING CHANGE: asdf`

	footers := MatchFooters(testStr)
	if len(footers) != 4 {
		t.Errorf("MatchFooters failed, empty input expected matched %d, matched %d, %v", 4, len(footers), footers)
	}

	var tests = []TestStrInOut{
		{footers[0], "Acked-By: RT fix readme" + NL + "with second line"},
		{footers[1], "Reviewed-By: RT"},
		{footers[2], "fix #1"},
		{footers[3], "BREAKING CHANGE: asdf"},
	}

	for _, test := range tests {
		if got, expected := test.in, test.out; got != expected {
			t.Errorf("MatchFooters failed, expected: %q, got: %q", expected, got)
		}
	}
}

func TestParseFooter(t *testing.T) {
	if token, sep, val := ParseFooter("token: "); token != "token" || sep != FSepColonSpace || val != "" {
		t.Error("ParseFooter check failed")
	}

	// footer's value can have newlines
	if token, sep, val := ParseFooter(fmt.Sprintf("token: value1%svalue2", NL)); token != "token" || sep != FSepColonSpace || val != "value1"+NL+"value2" {
		t.Error("ParseFooter check failed")
	}

	if token, sep, val := ParseFooter("fix #1"); token != "fix" || sep != FSepSpaceSharp || val != "1" {
		t.Error("ParseFooter check failed")
	}

	if token, sep, val := ParseFooter("fix sth"); token != "" || sep != "" || val != "" {
		t.Error("ParseFooter check failed")
	}
}

func TestTrimFooter(t *testing.T) {
	var tests = []TestStrInOut{
		{" re #1 ", "re #1"},
		{"Reviewed-by: some author " + NL, "Reviewed-by: some author"},
		{"Acked-by: ", "Acked-by: "},
		{"BREAKING CHANGE: ", "BREAKING CHANGE: "},
		// preserver separators
		{"  #1 ", " #1"},
		{" Reviewed-by:   ", "Reviewed-by: "},
	}

	for _, test := range tests {
		if got, expected := TrimFooter(test.in), test.out; got != expected {
			t.Errorf("TrimFooter failed, expected: %q, got: %q", expected, got)
		}
	}
}

// TestCommitMsgFooter
// @spec conventional commits v1.0.0
// 8.
//  a. (template) One or more footers MAY be provided one blank line after the body
//  b. (validate) Each footer MUST consist of a word token, followed by either a :<space> or <space># separator, followed by a string
// 9. (validate) A footer’s token MUST use `-` in place of whitespace characters
// 10. footer’s value MAY contain spaces and newlines, and parsing MUST terminate when the next valid footer token/separator pair is observed
// 12. if included as a footer, a breaking change MUST consist of the uppercase text BREAKING CHANGE, followed by a colon, space, and description
// 16. BREAKING-CHANGE MUST be synonymous with BREAKING CHANGE, when used as a token in a footer.
func TestCommitMsgFooter(t *testing.T) {
	// test validation
	validFts := []string{
		fmt.Sprintf("%s: some%schange%sof lines", FTokenBrkChange, NL, NL),
		fmt.Sprintf("%s: some%schange%sof lines", FTokenBrkChangeAlias, NL, NL),
		"Acked-by: RT",
		"Reviewed: ",
		"fix #1",
	}

	for _, f := range validFts {
		validCmtMsg := makeCommitMsg("fix", "spec 8b", false, "test footer", "", []string{f})
		if ok, msg := validCmtMsg.Validate(); !ok {
			t.Errorf("Spec rule 8b check failed, footer: %q should be valid, got msg: %q", f, msg)
		}
	}

	invalidFts := []string{
		": some change",        // no token
		" #1",                  // no token
		"footer some change",   // no separator
		"token 2: some change", // whitespace in token
		"token	2: some change", // whitespace in token \t
		"token\n2: some change",                         // whitespace in token \n
		"token\r\n2: some change",                       // whitespace in token \r\n
		fmt.Sprintf("%s: ", FTokenBrkChange),            // breaking change description is required if included in footer
		fmt.Sprintf("%s: ", FTokenBrkChangeAlias),       // breaking change description is required if included in footer
		fmt.Sprintf("%s #some change", FTokenBrkChange), // breaking change should use colon space separator
	}

	for _, f := range invalidFts {
		invalidCmtMsg := makeCommitMsg("fix", "spec 8b", false, "test footer", "", []string{f})
		if ok, msg := invalidCmtMsg.Validate(); ok {
			t.Errorf("Spec rule 8b check failed, commit footer: %q should be invalid with msg: %q", f, msg)
		}
	}

	footerAfterBodyMsg := makeCommitMsg("test", "spec 8a", false, "check newline after body", "body", []string{"Acked-by: RT"})
	if got, expected := footerAfterBodyMsg.ToString(), strings.ReplaceAll("test(spec 8a): check newline after body%s%sbody%s%sAcked-by: RT%s", "%s", NL); got != expected {
		t.Errorf("Spec rule 8a check failed, expected: %q, got: %q", expected, got)
	}

	footerAfterHeaderMsg := makeCommitMsg("test", "spec 8a", false, "check newline after header", "", []string{"Acked-by: RT"})
	if got, expected := footerAfterHeaderMsg.ToString(), strings.ReplaceAll("test(spec 8a): check newline after header%s%sAcked-by: RT%s", "%s", NL); got != expected {
		t.Errorf("Spec rule 8a check failed, expected: %q, got: %q", expected, got)
	}
}

func TestCommitCmdRun(t *testing.T) {
	// test flags mode functions
	if err := commitCmd.Flags().Set("type", "fix"); err != nil {
		t.Error("get string flag error, error expected: nil, got: ", err)
	}

	if cmtType, err := commitCmd.LocalFlags().GetString("type"); err != nil {
		t.Error("get string flag error, error expected: nil, got: ", err)
	} else if cmtType != "fix" {
		t.Errorf("get string flag error, type expected: %q, got: %q", "fix", cmtType)
	}

	// test for os.exit(1)
	// if os.Getenv("ShouldCommitCmdRunCrash") == "1" {
	// 	commitCmd.Run(commitCmd, []string{})
	// 	return
	// }
	// cmd := exec.Command(os.Args[0], "-test.run=TestCommitCmdRun")
	// cmd.Env = append(os.Environ(), "ShouldCommitCmdRunCrash=1")
	// err := cmd.Run()
	// if e, ok := err.(*exec.ExitError); ok && !e.Success() {
	// 	return
	// }
	// t.Fatalf("TestCommitCmdRun ran with err %v, want exit status 1", err)
}
