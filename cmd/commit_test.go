package cmd

import (
	"fmt"
	"runtime"
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

// TestMakeCommitMsg test trim string input
// - whitespaces including `\r\n` on Windows should be trimmed
func TestMakeCommitMsg(t *testing.T) {
	msg := makeCommitMsg(" fix ", " lib ", false, " desc\r\n", " bodyln1\n bodyln2\n", []string{" Acked-by: RT \n", "Review-by: RT2\r\n"})
	if msg.Type != "fix" ||
		msg.Scope != "lib" ||
		msg.Description != "desc" ||
		msg.Body != "bodyln1\n bodyln2" ||
		msg.Footers[0] != "Acked-by: RT" ||
		msg.Footers[1] != "Review-by: RT2" {
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
	msg1 := makeCommitMsg("docs", "", false, "fix typo", "msg body", []string{})

	if s := msg1.ToString(TestCmtMsgTmplPath); s != fmt.Sprintf("docs: fix typo%s%smsg body%s", NL, NL, NL) {
		t.Errorf(`expected: %s, got: %s`, `docs: fix typo\n\nmsg body\n`, s)
	}
}

func TestCommitMsgFooter(t *testing.T) {
}
