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
	if msg.Type != "fix" ||
		msg.Scope != "lib" ||
		msg.Description != "desc" ||
		msg.Body != "bodyln1\n bodyln2" ||
		msg.Footers[0] != "Acked-by: RT" ||
		msg.Footers[1] != "Review-by: RT2" ||
		msg.Footers[2] != " #1" ||
		msg.Footers[3] != "Reviewed: " {
		t.Error("Commit msg string input trim error")
	}
}

// TestCommitMsgHeader
// @spec conventional commits v1.0.0
// 1. REQUIRED `type` of a noun, OPTIONAL `scope`, OPTIONAL `!``, and REQUIRED terminal colon and space.
// 4. `scope` MUST consist of a noun describing a section of the codebase surrounded by parenthesis
// 5. `description` MUST immediately follow the colon and space after the type/scope prefix.
func TestCommitMsgHeader(t *testing.T) {
	// test validation
	// - required `type` and `description`
	// - no linebreaks within `type`, `scope`, and `description`
	emptyMsg := makeCommitMsg("", "", false, "", "", []string{})
	invalidScopeMsg := makeCommitMsg("fix", "with\nlinebreak", false, "", "", []string{})
	noDescMsg := makeCommitMsg("fix", "", false, "", "", []string{})

	if ok, msg := emptyMsg.Validate(); ok || msg != RequiredType {
		t.Errorf("Required commit type check failed, expected: %s, got: %s", RequiredType, msg)
	}

	if ok, msg := invalidScopeMsg.Validate(); ok || msg != InvalidScope {
		t.Errorf("No linebreak in scope check failed, expected: %s, got: %s", InvalidScope, msg)
	}

	if ok, msg := noDescMsg.Validate(); ok || msg != RequiredDesc {
		t.Errorf("Required commit description check failed, expected: %s, got: %s", RequiredDesc, msg)
	}

	// test toString format
	msg1 := makeCommitMsg("docs", "", false, "fix typo", "", []string{})           // test type trim
	msg2 := makeCommitMsg("docs", "READ ME.md", false, "fix typo", "", []string{}) // test scope trim
	msg3 := makeCommitMsg("docs", "", true, "fix typo", "", []string{})
	msg4 := makeCommitMsg("fix", "lib", true, "fix bug", "", []string{})

	if got, expected := msg1.ToString(), fmt.Sprintln("docs: fix typo"); got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}

	if got, expected := msg2.ToString(), fmt.Sprintln("docs(READ ME.md): fix typo"); got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}

	if got, expected := msg3.ToString(), fmt.Sprintln("docs!: fix typo"); got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}

	if got, expected := msg4.ToString(), fmt.Sprintln("fix(lib)!: fix bug"); got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}
}

// TestCommitMsgBody
// @spec conventional commits v1.0.0
// 6. body MUST begin one blank line after the description
func TestCommitMsgBody(t *testing.T) {
	body := fmt.Sprintln("msg body") + NL + "body line2"
	msg1 := makeCommitMsg("docs", "", false, "fix typo", body, []string{})

	if got, expected := msg1.ToString(), fmt.Sprintln("docs: fix typo")+NL+fmt.Sprintln(body); got != expected {
		t.Errorf("expected: %s%%, got: %s%%", expected, got)
	}
}

func TestParseFooter(t *testing.T) {
	if token, sep, val := ParseFooter("token: "); token != "token" || sep != FSepColonSpace || val != "" {
		t.Error("ParseFooter check failed")
	}

	// footer's value can have newlines
	if token, sep, val := ParseFooter(fmt.Sprintf("token: value1%svalue2", NL)); token != "token" || sep != FSepColonSpace || val != fmt.Sprintf("value1%svalue2", NL) {
		t.Error("ParseFooter check failed")
	}

	if token, sep, val := ParseFooter("fix #1"); token != "fix" || sep != FSepSpaceSharp || val != "1" {
		t.Error("ParseFooter check failed")
	}

	if token, sep, val := ParseFooter("fix sth"); token != "" || sep != "" || val != "" {
		t.Error("ParseFooter check failed")
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
		"Reviewed: ", // separator space should not be trimed if no footer value
		"fix #1",
	}

	for _, f := range validFts {
		validCmtMsg := makeCommitMsg("fix", "spec 8b", false, "test footer", "", []string{f})
		if ok, msg := validCmtMsg.Validate(); !ok {
			t.Errorf("Spec rule 8b check failed, footer: %s should be valid, got msg: %s", f, msg)
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
			t.Errorf("Spec rule 8b check failed, commit footer: %s should be invalid with msg: %s", f, msg)
		}
	}

	footerAfterBodyMsg := makeCommitMsg("test", "spec 8a", false, "check newline after body", "body", []string{"Acked-by: RT"})
	if got, expected := footerAfterBodyMsg.ToString(), strings.ReplaceAll("test(spec 8a): check newline after body%s%sbody%s%sAcked-by: RT%s", "%s", NL); got != expected {
		t.Errorf(`Spec rule 8a check failed, expected: %s%%, got: %s%%`, expected, got)
	}

	footerAfterHeaderMsg := makeCommitMsg("test", "spec 8a", false, "check newline after header", "", []string{"Acked-by: RT"})
	if got, expected := footerAfterHeaderMsg.ToString(), strings.ReplaceAll("test(spec 8a): check newline after header%s%sAcked-by: RT%s", "%s", NL); got != expected {
		t.Errorf(`Spec rule 8a check failed, expected: %s%%, got: %s%%`, expected, got)
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
		t.Errorf("get string flag error, type expected: %s, got: %s", "fix", cmtType)
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
