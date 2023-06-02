package pkg

import (
	"encoding/json"
	"errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	outputFormat := OutputFormatType(viper.GetString("output"))

	switch outputFormat {
	case OutputYaml:
		yamlData, err := yaml.Marshal(&data)
		if err != nil {
			return err
		}
		pterm.Println(string(yamlData))
	case OutputJson:
		jsonData, err := json.MarshalIndent(&data, "", "  ")
		if err != nil {
			return err
		}
		pterm.Println(string(jsonData))
	case OutputPretty:
		if err := pterm.DefaultTable.WithHasHeader().WithData(prettyPrintFunction(data)).Render(); err != nil {
			return err
		}
	default:
		if err := pterm.DefaultTable.WithHasHeader().WithData(prettyPrintFunction(data)).Render(); err != nil {
			return err
		}
	}
	return nil
}
