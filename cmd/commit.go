package cmd

import (
	"bytes"
	"fmt"
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

// Validate commit msg elements
func (cm CommitMsg) Validate() bool {
	return true
}

// String format commit msg as conventional commits spec v1.0.0
func (cm CommitMsg) String() string {
	var tmplBytes bytes.Buffer

	tmpl := mustTmpl(template.ParseFiles("templates/commitmsg.tmpl")) // filepath relative to main.go
	must(tmpl.Execute(&tmplBytes, cm))

	return tmplBytes.String()
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "conventional commit",
	Long:  "pass no flag to use interactive mode",
	Run: func(cmd *cobra.Command, args []string) {
		// try construct commit msg from flags
		var cmtMsg CommitMsg
		// hasAnyFlag := false

		cmtMsg.Type = mustStr(cmd.Flags().GetString("type"))
		cmtMsg.Scope = mustStr(cmd.Flags().GetString("scope"))
		cmtMsg.HasBrkChange = mustBool(cmd.Flags().GetBool("breaking"))
		cmtMsg.Description = mustStr(cmd.Flags().GetString("description"))
		cmtMsg.Body = mustStr(cmd.Flags().GetString("body"))
		cmtMsg.Footers = mustStrSlice(cmd.Flags().GetStringSlice("footers"))

		fmt.Print(cmtMsg.String())
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
