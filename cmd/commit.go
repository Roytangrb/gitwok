package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// RequiredType error msg of missing type
	RequiredType = "commit type is required"
	// RequiredDesc error msg of missing description
	RequiredDesc = "commit description is required"
	// InvalidType error msg of invalid type
	InvalidType = "commit type is invalid"
	// InvalidScope error msg of invalid scope
	InvalidScope = "commit scope is invalid"
	// InvalidDesc error msg of invalid description
	InvalidDesc = "commit description is invalid"
	// InvalidBrkChnFTSep error msg of invalid breaking change footer separator
	InvalidBrkChnFTSep = "breaking change footer separator is invalid"
	// RequiredBrkChnFTDesc error msg of missing breaking footer description
	RequiredBrkChnFTDesc = "breaking change footer description is required"

	// InvalidFooter error msg of invalid footer
	InvalidFooter = "commit footer is invalid"
	// InvalidFooterToken error msg of invalid footer
	InvalidFooterToken = "commit footer token is invalid"
	// FTokenBrkChange special footer token
	FTokenBrkChange = "BREAKING CHANGE"
	// FTokenBrkChangeAlias FTokenBrkChange alias
	FTokenBrkChangeAlias = "BREAKING-CHANGE"
	// FSepColonSpace footer separator
	FSepColonSpace = ": "
	// FSepSpaceSharp footer separator
	FSepSpaceSharp = " #"
)

// CommitMsg properties
type CommitMsg struct {
	Type         string   `survey:"type"`        // required, preset or config values only
	Scope        string   `survey:"scope"`       // optional
	HasBrkChange bool     `survey:"breaking"`    // optional, default false
	Description  string   `survey:"description"` // required, no line break
	Body         string   `survey:"body"`        // optional, allow line breaks
	Footers      []string `survey:"footers"`     // optional, allow multiple lines
}

// CommitMsgTmpl template for building commit message
const CommitMsgTmpl = `{{.Type}}{{if .Scope}}({{.Scope}}){{end}}{{if .HasBrkChange}}!{{end}}: {{.Description}}
{{if .Body}}
{{.Body}}
{{end}}{{if .Footers}}
{{- range .Footers}}
{{. -}}
{{end}}
{{end}}`

// PresetCommitTypes conventional commits suggested types
var PresetCommitTypes = []string{"fix", "feat", "build", "chore", "ci", "docs", "perf", "refactor", "style", "test"}

// CommitMsgQuestions build CommitMsg struct for interactive commit mode
var CommitMsgQuestions = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Choose commit type:",
			Options: PresetCommitTypes,
			Default: PresetCommitTypes[0],
		},
	},
	{
		Name: "scope", // TODO: use preset options from config
		Prompt: &survey.Input{
			Message: "Enter commit scope:",
		},
		Transform: survey.ComposeTransformers(
			// [Bug](https://github.com/AlecAivazis/survey/issues/329)
			survey.TransformString(strings.ToLower), // TODO: config option
			survey.TransformString(strings.TrimSpace),
		),
	},
	{
		Name: "breaking",
		Prompt: &survey.Confirm{
			Message: "Includes breaking changes?",
			Default: false,
		},
	},
	{
		Name: "description",
		Prompt: &survey.Input{
			Message: fmt.Sprintln("Enter commit description:"),
		},
		Validate:  survey.Required,
		Transform: survey.TransformString(strings.TrimSpace),
	},
	{
		Name: "body",
		Prompt: &survey.Multiline{
			Message: fmt.Sprintln("Enter optional detail description:"),
			// space is trim in survey multiline output answer
			// TODO: feat(survey) Transform API is missing in type Multiline
		},
	},
	// TODO: prompt footers and tranform to []string
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
		// separator space should not be trimed if no footer value
		if !strings.HasSuffix(s, FSepColonSpace) {
			s = strings.TrimRightFunc(s, unicode.IsSpace)
		}
		if !strings.HasPrefix(s, FSepSpaceSharp) {
			s = strings.TrimLeftFunc(s, unicode.IsSpace)
		}
		footers[i] = s
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

// ContainsNewline check if string contains newline chars
func ContainsNewline(s string) bool {
	return strings.Contains(s, fmt.Sprintln())
}

// ContainsWhiteSpace check if string contains whitesapces
func ContainsWhiteSpace(s string) bool {
	for _, c := range s {
		if unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

// IsBrkChnFooter check if token is breaking change footer token
func IsBrkChnFooter(token string) bool {
	return token == FTokenBrkChange || token == FTokenBrkChangeAlias
}

// ParseFooter return components of a commit msg footer if seperable by ": " or " #"
// @param f footer without no newlines
// @return token "" if separated wrongly
// @return sep "" if separated wrongly
// @return val "" if separated wrongly
func ParseFooter(f string) (token, sep, val string) {
	if elms := strings.Split(f, FSepColonSpace); len(elms) == 2 {
		token, sep, val = elms[0], FSepColonSpace, elms[1]
		return
	} else if elms := strings.Split(f, FSepSpaceSharp); len(elms) == 2 {
		token, sep, val = elms[0], FSepSpaceSharp, elms[1]
		return
	}

	return "", "", ""
}

// Validate commit msg elements
// @return valid {bool}
// @return msg {string} error msg
func (cm *CommitMsg) Validate() (bool, string) {
	if cm.Type == "" {
		return false, RequiredType
	} else if ContainsNewline(cm.Type) {
		return false, InvalidType
	}

	if cm.Scope != "" && ContainsNewline(cm.Scope) {
		return false, InvalidScope
	}

	if cm.Description == "" {
		return false, RequiredDesc
	} else if ContainsNewline(cm.Description) {
		return false, InvalidDesc
	}

	if len(cm.Footers) > 0 {
		for _, f := range cm.Footers {
			token, sep, val := ParseFooter(f)
			if token == "" || sep == "" {
				return false, InvalidFooter
			}
			if token != FTokenBrkChange && ContainsWhiteSpace(token) {
				return false, InvalidFooterToken
			}
			if IsBrkChnFooter(token) && sep != FSepColonSpace {
				return false, InvalidBrkChnFTSep
			}
			if IsBrkChnFooter(token) && val == "" {
				return false, RequiredBrkChnFTDesc
			}
		}
	}

	return true, ""
}

// ToString format commit msg as conventional commits spec v1.0.0
func (cm *CommitMsg) ToString() string {
	var tmplBytes bytes.Buffer

	tmpl := mustTmpl(template.New("commitmsg").Parse(CommitMsgTmpl))
	must(tmpl.Execute(&tmplBytes, cm))

	return tmplBytes.String()
}

// Commit validate and git commit the CommitMsg
func (cm *CommitMsg) Commit() {
	if ok, msg := cm.Validate(); ok {
		cmtMsgStr := cm.ToString()
		if isVerbose {
			logger.Verbose("Executing git commit -m with msg: ")
			fmt.Print(cmtMsgStr)
		}
		GitCommit(cmtMsgStr)
	} else {
		logger.Fatal(msg)
	}
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "build and make conventional commit",
	Long:  "Pass no flag to use interactive mode or build commit message with flags",
	Run: func(cmd *cobra.Command, args []string) {
		// count local flags set explicitly
		// cobra issue: https://github.com/spf13/cobra/issues/1315
		flagCount := 0
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				flagCount++
			}
		})

		// use flags mode if any flag has been set
		if flagCount > 0 {
			// readonly local flags
			cmtType := mustStr(cmd.LocalFlags().GetString("type"))
			cmtScope := mustStr(cmd.LocalFlags().GetString("scope"))
			cmtHasBrkChange := mustBool(cmd.LocalFlags().GetBool("breaking"))
			cmtDescription := mustStr(cmd.LocalFlags().GetString("description"))
			cmtBody := mustStr(cmd.LocalFlags().GetString("body"))
			cmtFooters := mustStrSlice(cmd.LocalFlags().GetStringSlice("footers"))

			// try construct commit msg from flags
			cmtMsg := makeCommitMsg(cmtType, cmtScope, cmtHasBrkChange, cmtDescription, cmtBody, cmtFooters)
			cmtMsg.Commit()
		} else {
			var cmtMsg CommitMsg
			must(survey.Ask(CommitMsgQuestions, &cmtMsg))
			cmtMsg.Commit()
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("type", "t", "", "required: commit type")
	commitCmd.Flags().StringP("scope", "s", "", "optional: commit scope")
	commitCmd.Flags().BoolP("breaking", "k", false, "optional: has breaking change")
	commitCmd.Flags().StringP("description", "d", "", "required: commit description")
	commitCmd.Flags().StringP("body", "b", "", "optional: commit body")
	commitCmd.Flags().StringSliceP("footers", "f", []string{}, "optional: commit footers, allow multiple")
}

func mustTmpl(tmpl *template.Template, err error) *template.Template {
	must(err)
	return tmpl
}
