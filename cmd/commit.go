package cmd

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

var cmtScopeInput = &survey.Question{
	Name: "scope",
	Prompt: &survey.Input{
		Message: "Enter commit scope:",
	},
	Transform: survey.ComposeTransformers(
		// TODO: [Bug](https://github.com/AlecAivazis/survey/issues/329)
		// survey.TransformString(strings.ToLower), // TODO: config option
		survey.TransformString(strings.TrimSpace),
	),
}

var cmtBrkConfirm = &survey.Question{
	Name: "breaking",
	Prompt: &survey.Confirm{
		Message: "Includes breaking changes?",
		Default: false,
	},
}

var cmtDescInput = &survey.Question{
	Name: "description",
	Prompt: &survey.Input{
		Message: "Enter commit description:",
	},
	Validate:  survey.Required,
	Transform: survey.TransformString(strings.TrimSpace),
}

var cmtBodyMulti = &survey.Question{
	Name: "body",
	Prompt: &survey.Multiline{
		Message: "Enter optional commit body:",
		// space is trim in survey multiline output answer
		// TODO: feat(survey) Transform API is missing in type Multiline
	},
}

// FootersQuestions ask footers input separately to handle custom parse logic
var FootersQuestions = []*survey.Question{
	{
		Name: "footers",
		Prompt: &survey.Multiline{
			Message: "Enter optional footers:",
		},
	},
}

// CommitFooters helper struct for separate survey with custom Setter
type CommitFooters struct {
	Footers []string `survey:"footers"`
}

// WriteAnswer implements Settable interface of CommitMsg for survey
// assign string input to parsed footer slice
func (ft *CommitFooters) WriteAnswer(name string, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("Write %s error, got: %v", name, value)
	}

	if str == "" {
		ft.Footers = []string{}
		return nil
	}

	ft.Footers = MatchFooters(str)

	return nil
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
		footers[i] = TrimFooter(s)
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

// MatchFooters takes raw footers input string and return
// string slice of found footers.
// Match <token + sep>, find indices and take values in between
// Caveats: "BREAKING CHANGE #" is not matched, for " #" is not
// valid separator for breaking change, instead, "CHANGE #" is matched
func MatchFooters(str string) []string {
	re := regexp.MustCompile(`([\w-]+(: | #))|(BREAKING CHANGE: )`)
	indices := re.FindAllStringIndex(str, -1)

	if indices == nil {
		logger.Warn("No valid footer message found")
		return []string{}
	}

	footers := []string{}
	for i, idx := range indices {
		start, end := idx[0], idx[1]
		tokenSep := str[start:end]
		var value string
		if i != len(indices)-1 {
			nextStart := indices[i+1][0]
			value = str[end:nextStart]
		} else {
			value = str[end:]
		}
		footer := TrimFooter(tokenSep + value)
		footers = append(footers, footer)
	}
	logger.Verbose(fmt.Sprintf("Parsed %d footers:", len(footers)), footers)

	return footers
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

// TrimFooter separator space should not be trimed if no footer value
// Whitespace trimmed value may contains leading "#" or trailing ":",
// treat it as intention to mean a separator, add the space back
func TrimFooter(s string) string {
	if !strings.HasSuffix(s, FSepColonSpace) {
		s = strings.TrimRightFunc(s, unicode.IsSpace)
	}
	if !strings.HasPrefix(s, FSepSpaceSharp) {
		s = strings.TrimLeftFunc(s, unicode.IsSpace)
	}

	if strings.HasPrefix(s, "#") {
		s = " " + s
	}

	if strings.HasSuffix(s, ":") {
		s = s + " "
	}

	return s
}

// Validate commit msg elements
// @return valid {bool}
// @return msg {string} error msg
func (cm *CommitMsg) Validate() (bool, string) {
	if cm.Type == "" {
		return false, RequiredType
	} else if ContainsWhiteSpace(cm.Type) {
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
func (cm *CommitMsg) Commit(git *Git) {
	if ok, msg := cm.Validate(); ok {
		cmtMsgStr := cm.ToString()
		logger.Verbose(fmt.Sprintln("Executing git commit -m with msg: ") + cmtMsgStr)

		git.Commit("-m", cmtMsgStr)
	} else {
		logger.Fatal(msg)
	}
}

// Prompt use interactive prompts to build the commit message
func (cm *CommitMsg) Prompt() {
	var questions = []*survey.Question{}

	// prompt type
	var typeOptions []string
	if options := viper.GetStringSlice("gitwok.commit.type"); len(options) != 0 {
		typeOptions = options
	} else {
		typeOptions = PresetCommitTypes
	}

	questions = append(questions, &survey.Question{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Choose commit type:",
			Options: typeOptions,
			Default: typeOptions[0],
		},
	})

	// prompt scope
	if prompt := viper.GetBool("gitwok.commit.prompt.scope"); prompt {
		if options := viper.GetStringSlice("gitwok.commit.scope"); len(options) != 0 {
			var cmtScopeSelect = &survey.Question{
				Name: "scope",
				Prompt: &survey.Select{
					Message: "Choose commit scope:",
					Options: options,
					Default: options[0],
				},
			}
			questions = append(questions, cmtScopeSelect)
		} else {
			questions = append(questions, cmtScopeInput)
		}
	}

	// prompt breaking
	if prompt := viper.GetBool("gitwok.commit.prompt.breaking"); prompt {
		questions = append(questions, cmtBrkConfirm)
	}

	// prompt description
	questions = append(questions, cmtDescInput)

	// prompt body
	if prompt := viper.GetBool("gitwok.commit.prompt.body"); prompt {
		questions = append(questions, cmtBodyMulti)
	}

	must(survey.Ask(questions, cm))

	// prompt footers
	if prompt := viper.GetBool("gitwok.commit.prompt.footers"); prompt {
		var ft CommitFooters
		must(survey.Ask(FootersQuestions, &ft))
		cm.Footers = ft.Footers
	}
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "build and make conventional commit",
	Long:  "Pass no flag to use interactive mode or build commit message with flags",
	Run: func(cmd *cobra.Command, args []string) {
		var git = &Git{
			verbose: false,
			dryRun:  mustBool(cmd.Flags().GetBool("dry-run")),
		}

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
			cmtMsg.Commit(git)
		} else {
			var cmtMsg CommitMsg
			cmtMsg.Prompt()
			cmtMsg.Commit(git)
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
