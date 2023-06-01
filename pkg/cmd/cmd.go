package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"gobit/pkg/cmd/config"
	"gobit/pkg/cmd/permission"
	"gobit/pkg/cmd/project"
	"gobit/pkg/cmd/repository"
	"gobit/pkg/cmd/version"
	"gobit/pkg/cmd/webhook"
	"os"
	"path/filepath"
)

var (
	cfgFile      string
	baseUrl      string
	userToken    string
	limit        int
	outputFormat pkg.OutputFormatType

	rootCmd = &cobra.Command{
		Use:   "gobit",
		Short: "GoBit - Simple CLI-tool for Bitbucket",
		Long:  `GoBit – A Simple CLI-Tool for Bitbucket written in Go using Cobra`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default $HOME/.gobit/config.yaml")
	rootCmd.PersistentFlags().StringVar(&baseUrl, "baseUrl", "https://git.spk.no", "Base url for BitBucket instance")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 100, "Max return values")
	rootCmd.PersistentFlags().StringVarP(&userToken, "token", "t", "", "Token for user")
	rootCmd.PersistentFlags().VarP(&outputFormat, "output", "o", "Output format. One of: pretty, yaml, json")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("baseUrl", rootCmd.PersistentFlags().Lookup("baseUrl"))
	viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(permission.Cmd)
	rootCmd.AddCommand(project.Cmd)
	rootCmd.AddCommand(repository.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(webhook.Cmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, ".gobit"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		cfgFile = filepath.Join(home, ".gobit", "config.yaml")
	}

	viper.AutomaticEnv()

	pkg.CreateFileIfNotExists(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		pterm.Error.Println("Error reading config file:", viper.ConfigFileUsed())
	}

	viper.SetDefault("config", viper.ConfigFileUsed())
}
