package logger

import (
	"errors"
	"fmt"
	"git.spk.no/infra/bucketctl/pkg/common"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

type LogFormatType string

const (
	LogFormatPretty LogFormatType = "pretty"
	LogFormatPlain  LogFormatType = "plain"
)

func LogFormatCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"pretty\tOutput in pretty format",
		"plain\tOutput in plain text",
	}, cobra.ShellCompDirectiveDefault
}

func (e *LogFormatType) String() string {
	return string(*e)
}

func (e *LogFormatType) Set(v string) error {
	switch v {
	case "pretty", "plain":
		*e = LogFormatType(v)
		return nil
	default:
		return errors.New(`must be one of "pretty" or "plain"`)
	}
}

func (e *LogFormatType) Type() string {
	return "LogFormatType"
}

func Log(str string, args ...interface{}) {
	switch LogFormatType(viper.GetString(common.LogFormatFlag)) {
	case LogFormatPlain:
		pterm.Printfln(plainSprintf(str, args...))
	case LogFormatPretty:
		pterm.Printfln(str, args...)
	}
}

func Info(str string, args ...interface{}) {
	switch LogFormatType(viper.GetString(common.LogFormatFlag)) {
	case LogFormatPlain:
		pterm.Info.Printfln(plainSprintf(str, args...))
	case LogFormatPretty:
		pterm.Info.Printfln(str, args...)
	}
}

func Warn(str string, args ...interface{}) {
	switch LogFormatType(viper.GetString(common.LogFormatFlag)) {
	case LogFormatPlain:
		pterm.Warning.Printfln(plainSprintf(str, args...))
	case LogFormatPretty:
		pterm.Warning.Printfln(str, args...)
	}
}

func Err(str string, args ...interface{}) {
	switch LogFormatType(viper.GetString(common.LogFormatFlag)) {
	case LogFormatPlain:
		pterm.Error.Printfln(plainSprintf(str, args...))
	case LogFormatPretty:
		pterm.Error.Printfln(str, args...)
	}
}

var nonPrintableRegEx = regexp.MustCompile(`[^[:ascii:]\p{Latin}]`)

// plainSprintf takes a "pretty" string and strips non-ascii, non-latin characters (colours, emojis, etc.),
// returning a "plain" string
//
// Parameters:
//   - pretty: The format string to be processed.
//   - args: Optional list of arguments to be formatted into the string.
//
// Returns:
//   - plain: A plain string containing only allowed characters
func plainSprintf(pretty string, args ...interface{}) (plain string) {
	plain = pterm.RemoveColorFromString(fmt.Sprintf(pretty, args...))
	plain = nonPrintableRegEx.ReplaceAllString(plain, "")
	plain = strings.TrimSpace(plain)
	pterm.DisableColor()
	return plain
}
