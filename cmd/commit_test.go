package cmd

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

const TestCmtMsgTmplPath = "../templates/commitmsg.tmpl"

// NL is newline represented in os file
var NL string = DefineNewline()

func DefineNewline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func TestContainsNewline(t *testing.T) {
	if !ContainsNewline(NL) {
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
		t.Errorf(`Required commit type check failed, expected: %s, got: %s`, RequiredType, msg)
	}

	if ok, msg := invalidScopeMsg.Validate(); ok || msg != InvalidScope {
		t.Errorf(`No linebreak in scope check failed, expected: %s, got: %s`, InvalidScope, msg)
	}

	if ok, msg := noDescMsg.Validate(); ok || msg != RequiredDesc {
		t.Errorf(`Required commit description check failed, expected: %s, got: %s`, RequiredDesc, msg)
	}

	// test toString format
	msg1 := makeCommitMsg("docs", "", false, "fix typo", "", []string{})           // test type trim
	msg2 := makeCommitMsg("docs", "READ ME.md", false, "fix typo", "", []string{}) // test scope trim
	msg3 := makeCommitMsg("docs", "", true, "fix typo", "", []string{})
	msg4 := makeCommitMsg("fix", "lib", true, "fix bug", "", []string{})

	if s := msg1.ToString(TestCmtMsgTmplPath); s != fmt.Sprintf("docs: fix typo%s", NL) {
		t.Errorf(`expected: %s, got: %s`, `docs: fix typo\n`, s)
	}

	if s := msg2.ToString(TestCmtMsgTmplPath); s != fmt.Sprintf("docs(READ ME.md): fix typo%s", NL) {
		t.Errorf(`expected: %s, got: %s`, `docs(READ ME.md): fix typo\n`, s)
	}

	if s := msg3.ToString(TestCmtMsgTmplPath); s != fmt.Sprintf("docs!: fix typo%s", NL) {
		t.Errorf(`expected: %s, got: %s`, `docs!: fix typo\n`, s)
	}

	if s := msg4.ToString(TestCmtMsgTmplPath); s != fmt.Sprintf("fix(lib)!: fix bug%s", NL) {
		t.Errorf(`expected: %s, got: %s`, `fix(lib)!: fix bug\n`, s)
	}
}

// TestCommitMsgBody
// @spec conventional commits v1.0.0
// 6. body MUST begin one blank line after the description
func TestCommitMsgBody(t *testing.T) {
	msg1 := makeCommitMsg("docs", "", false, "fix typo", "msg body\nbody line2", []string{})

	if s := msg1.ToString(TestCmtMsgTmplPath); s != fmt.Sprintf("docs: fix typo%s%smsg body\nbody line2%s", NL, NL, NL) {
		t.Errorf(`expected: %s, got: %s`, `docs: fix typo\n\nmsg body\nbody line2\n`, s)
	}
}

func TestParseFooter(t *testing.T) {
	if token, sep, val := ParseFooter("token: "); token != "token" || sep != FSepColonSpace || val != "" {
		t.Error("ParseFooter check failed")
	}

	// footer's value can have newlines
	if token, sep, val := ParseFooter("token: value1\nvalue2"); token != "token" || sep != FSepColonSpace || val != "value1\nvalue2" {
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
		fmt.Sprintf("%s: some\nchange\nof lines", FTokenBrkChange),
		fmt.Sprintf("%s: some\nchange\nof lines", FTokenBrkChangeAlias),
		"Acked-by: RT",
		"Reviewed: ", // seperator space should not be trimed if no footer value
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
		"footer some change",   // no seperator
		"token 2: some change", // whitespace in token
		"token	2: some change", // whitespace in token \t
		"token\n2: some change",                         // whitespace in token \n
		"token\r\n2: some change",                       // whitespace in token \r\n
		fmt.Sprintf("%s: ", FTokenBrkChange),            // breaking change description is required if included in footer
		fmt.Sprintf("%s: ", FTokenBrkChangeAlias),       // breaking change description is required if included in footer
		fmt.Sprintf("%s #some change", FTokenBrkChange), // breaking change should use colon space seperator
	}

	for _, f := range invalidFts {
		invalidCmtMsg := makeCommitMsg("fix", "spec 8b", false, "test footer", "", []string{f})
		if ok, msg := invalidCmtMsg.Validate(); ok {
			t.Errorf("Spec rule 8b check failed, commit footer: %s should be invalid with msg: %s", f, msg)
		}
	}

	footerAfterBodyMsg := makeCommitMsg("test", "spec 8a", false, "check newline after body", "body", []string{"Acked-by: RT"})
	if s := footerAfterBodyMsg.ToString(TestCmtMsgTmplPath); s != strings.ReplaceAll("test(spec 8a): check newline after body%s%sbody%s%sAcked-by: RT%s", "%s", NL) {
		t.Errorf(`Spec rule 8a check failed, expected: %s, got: %s`, `test(spec 8a): check newline after body\n\nbody\n\nAcked-by: RT\n`, s)
	}

	footerAfterHeaderMsg := makeCommitMsg("test", "spec 8a", false, "check newline after header", "", []string{"Acked-by: RT"})
	if s := footerAfterHeaderMsg.ToString(TestCmtMsgTmplPath); s != strings.ReplaceAll("test(spec 8a): check newline after header%s%sAcked-by: RT%s", "%s", NL) {
		t.Errorf(`Spec rule 8a check failed, expected: %s, got: %s`, `test(spec 8a): check newline after header\n\nAcked-by: RT\n`, s)
	}
}
