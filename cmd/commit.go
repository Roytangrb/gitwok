package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
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
	return strings.Contains(s, "\n") || strings.Contains(s, "\r\n")
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
func (cm CommitMsg) Validate() (bool, string) {
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
func (cm CommitMsg) ToString(fp string) string {
	var tmplBytes bytes.Buffer

	tmpl := mustTmpl(template.ParseFiles(fp)) // filepath relative to main.go
	must(tmpl.Execute(&tmplBytes, cm))

	return tmplBytes.String()
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "build and make conventional commit",
	Long:  "Pass no flag to use interactive mode or build commit message with flags",
	Run: func(cmd *cobra.Command, args []string) {
		// use flags mode if any flag has been set
		flagMode := cmd.Flags().NFlag() > 0

		if flagMode {
			// readonly local flags
			cmtType := mustStr(cmd.LocalFlags().GetString("type"))
			cmtScope := mustStr(cmd.LocalFlags().GetString("scope"))
			cmtHasBrkChange := mustBool(cmd.LocalFlags().GetBool("breaking"))
			cmtDescription := mustStr(cmd.LocalFlags().GetString("description"))
			cmtBody := mustStr(cmd.LocalFlags().GetString("body"))
			cmtFooters := mustStrSlice(cmd.LocalFlags().GetStringSlice("footers"))

			// try construct commit msg from flags
			cmtMsg := makeCommitMsg(cmtType, cmtScope, cmtHasBrkChange, cmtDescription, cmtBody, cmtFooters)

			if ok, msg := cmtMsg.Validate(); ok {
				cmtMsgStr := cmtMsg.ToString("templates/commitmsg.tmpl")
				if isVerbose {
					logger.Verbose("Executing git commit -m with msg: ")
					fmt.Print(cmtMsgStr)
				}
				GitCommit(cmtMsgStr)
			} else {
				logger.Fatal(msg)
			}
		} else {
			// interactive mode
			fmt.Println("Using interactive mode")
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
