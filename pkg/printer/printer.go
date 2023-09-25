package printer

import (
	"bucketctl/pkg/common"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	ALL  = pterm.Underscore.Sprint("ALL")
	AUTH = pterm.Bold.Sprint("AUTHENTICATED")
)

type OutputFormatType string

const (
	OutputPretty OutputFormatType = "pretty"
	OutputYaml   OutputFormatType = "yaml"
	OutputJson   OutputFormatType = "json"
)

func OutputCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"pretty\tOutput in pretty format",
		"yaml\tOutput in YAML format",
		"json\tOutput in JSON format",
	}, cobra.ShellCompDirectiveDefault
}

func (e *OutputFormatType) String() string {
	return string(*e)
}

func (e *OutputFormatType) Set(v string) error {
	switch v {
	case "pretty", "yaml", "json":
		*e = OutputFormatType(v)
		return nil
	default:
		return errors.New(`must be one of "pretty", "yaml", or "json"`)
	}
}

func (e *OutputFormatType) Type() string {
	return "OutputFormatType"
}

func PrintData[T interface{}](data T, prettyPrintFunction func(a T) [][]string) error {
	outputFormat := OutputFormatType(viper.GetString(common.OutputFlag))

	if prettyPrintFunction == nil && (outputFormat == OutputPretty || outputFormat == "") {
		outputFormat = OutputYaml
	}

	switch outputFormat {
	case OutputYaml:
		var b bytes.Buffer
		yamlEncoder := yaml.NewEncoder(&b)
		yamlEncoder.SetIndent(2)

		if err := yamlEncoder.Encode(&data); err != nil {
			return err
		}
		pterm.Println(string(b.Bytes()))
	case OutputJson:
		jsonData, err := json.MarshalIndent(&data, "", "  ")
		if err != nil {
			return err
		}
		pterm.Println(string(jsonData))
	default:
		if err := pterm.DefaultTable.WithHasHeader().WithData(prettyPrintFunction(data)).Render(); err != nil {
			return err
		}
	}
	return nil
}
