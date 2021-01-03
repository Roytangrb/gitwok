package cmd

import "testing"

const TestCmtMsgTmplPath = "../templates/commitmsg.tmpl"

func TestCommitMsgHeader(t *testing.T) {
	msg1 := makeCommitMsg(" docs ", "", false, "fix typo", "", []string{})           // test type trim
	msg2 := makeCommitMsg("docs", " READ ME.md ", false, "fix typo", "", []string{}) // test scope trim
	msg3 := makeCommitMsg("docs", "", true, "fix typo", "", []string{})
	msg4 := makeCommitMsg("fix", "lib", true, "fix bug", "", []string{})

	if s := msg1.ToString(TestCmtMsgTmplPath); s != "docs: fix typo\n" {
		t.Errorf(`expected: %s, got: %s`, `docs: fix typo\n`, s)
	}

	if s := msg2.ToString(TestCmtMsgTmplPath); s != "docs(READ ME.md): fix typo\n" {
		t.Errorf(`expected: %s, got: %s`, `docs(README.md): fix typo\n`, s)
	}

	if s := msg3.ToString(TestCmtMsgTmplPath); s != "docs!: fix typo\n" {
		t.Errorf(`expected: %s, got: %s`, `docs!: fix typo\n`, s)
	}

	if s := msg4.ToString(TestCmtMsgTmplPath); s != "fix(lib)!: fix bug\n" {
		t.Errorf(`expected: %s, got: %s`, `fix(lib)!: fix bug\n`, s)
	}
}

func TestCommitMsgBody(t *testing.T) {
	msg1 := makeCommitMsg("docs", "", false, "fix typo", "msg body", []string{})

	if s := msg1.ToString(TestCmtMsgTmplPath); s != "docs: fix typo\n\nmsg body\n" {
		t.Errorf(`expected: %s, got: %s`, `docs: fix typo\n\nmsg body\n`, s)
	}
}

func TestCommitMsgFooter(t *testing.T) {
}

func TestCommitMsgLnBreaks(t *testing.T) {
}
