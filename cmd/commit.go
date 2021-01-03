package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// CommitMsg properties
type CommitMsg struct {
	Type         string   // required, preset or config values only
	Scope        string   // optional
	HasBrkChange bool     // optional, default false
	Description  string   // required, no line break
	Body         string   // optional, allow line breaks
	Footers      []string // optional, allow multiple lines
}

func makeCommitMsg(
	cmtType string,
	scope string,
	hasBrkChn bool,
	desc string,
	body string,
	footers []string,
) *CommitMsg {
	for i, s := range footers {
		footers[i] = strings.TrimSpace(s)
	}
	return &CommitMsg{
		Type:         strings.TrimSpace(cmtType),
		Scope:        strings.TrimSpace(scope),
		HasBrkChange: hasBrkChn,
		Description:  strings.TrimSpace(desc),
		Body:         strings.TrimSpace(body),
		Footers:      footers,
	}
}

// Validate commit msg elements
func (cm CommitMsg) Validate() bool {
	return true
}

// ToString format commit msg as conventional commits spec v1.0.0
func (cm CommitMsg) ToString(fp string) string {
	var tmplBytes bytes.Buffer

	tmpl := mustTmpl(template.ParseFiles(fp)) // filepath relative to main.go
	must(tmpl.Execute(&tmplBytes, cm))

	return tmplBytes.String()
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "conventional commit",
	Long:  "pass no flag to use interactive mode",
	Run: func(cmd *cobra.Command, args []string) {
		cmtType := mustStr(cmd.Flags().GetString("type"))
		cmtScope := mustStr(cmd.Flags().GetString("scope"))
		cmtHasBrkChange := mustBool(cmd.Flags().GetBool("breaking"))
		cmtDescription := mustStr(cmd.Flags().GetString("description"))
		cmtBody := mustStr(cmd.Flags().GetString("body"))
		cmtFooters := mustStrSlice(cmd.Flags().GetStringSlice("footers"))

		// try construct commit msg from flags
		cmtMsg := makeCommitMsg(cmtType, cmtScope, cmtHasBrkChange, cmtDescription, cmtBody, cmtFooters)

		logger.Verbose("commit msg struct:", cmtMsg)
		fmt.Print(cmtMsg.ToString("templates/commitmsg.tmpl"))
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// type is required msg contructed with flags
	commitCmd.Flags().StringP("type", "t", "", "required: commit type")
	commitCmd.Flags().StringP("scope", "s", "", "optional: commit scope")
	commitCmd.Flags().BoolP("breaking", "k", false, "optional: has breaking change")
	commitCmd.Flags().StringP("description", "d", "", "required: commit description")
	commitCmd.Flags().StringP("body", "b", "", "optional: commit body")
	// TODO: handle footers passed as flags
	commitCmd.Flags().StringSliceP("footers", "f", []string{}, "optional: commit footers, allow multiple")
}

func mustTmpl(tmpl *template.Template, err error) *template.Template {
	must(err)
	return tmpl
}
