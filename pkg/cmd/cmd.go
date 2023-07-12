package cmd

import (
	"bucketctl/pkg"
	"bucketctl/pkg/cmd/config"
	"bucketctl/pkg/cmd/permission"
	"bucketctl/pkg/cmd/project"
	"bucketctl/pkg/cmd/repository"
	"bucketctl/pkg/cmd/settings"
	"bucketctl/pkg/cmd/version"
	"bucketctl/pkg/cmd/webhook"
	"bucketctl/pkg/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Use:   "bucketctl",
		Short: "bucketctl - Simple CLI-tool for Bitbucket",
		Long:  `bucketctl – A Simple CLI-Tool for Bitbucket written in Go using Cobra`,
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, types.ConfigFlag, "", "Config file (default $HOME/.bucketctl/config.yaml")
	rootCmd.PersistentFlags().StringVar(&baseUrl, types.BaseUrlFlag, "", "Base url for BitBucket instance")
	rootCmd.PersistentFlags().IntVarP(&limit, types.LimitFlag, "l", 500, "Max return values")
	rootCmd.PersistentFlags().StringVarP(&userToken, types.TokenFlag, "t", "", "Http access token")
	rootCmd.PersistentFlags().VarP(&outputFormat, types.OutputFlag, "o", "Output format. One of: pretty, yaml, json")

	viper.BindPFlag(types.ConfigFlag, rootCmd.PersistentFlags().Lookup(types.ConfigFlag))
	viper.BindPFlag(types.BaseUrlFlag, rootCmd.PersistentFlags().Lookup(types.BaseUrlFlag))
	viper.BindPFlag(types.LimitFlag, rootCmd.PersistentFlags().Lookup(types.LimitFlag))
	viper.BindPFlag(types.TokenFlag, rootCmd.PersistentFlags().Lookup(types.TokenFlag))
	viper.BindPFlag(types.OutputFlag, rootCmd.PersistentFlags().Lookup(types.OutputFlag))

	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(permission.Cmd)
	rootCmd.AddCommand(project.Cmd)
	rootCmd.AddCommand(repository.Cmd)
	rootCmd.AddCommand(settings.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(webhook.Cmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, ".config", "bucketctl"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		cfgFile = filepath.Join(home, ".config", "bucketctl", "config.yaml")
	}

	viper.AutomaticEnv()

	if err := pkg.CreateDirIfNotExists(cfgFile, 0700); err != nil {
		cobra.CheckErr(err)
	}

	if err := pkg.CreateFileIfNotExists(cfgFile, 0600); err != nil {
		cobra.CheckErr(err)
	}

	if err := pkg.CheckFilePermission(cfgFile, 0600); err != nil {
		cobra.CheckErr(err)
	}

	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}

	viper.SetDefault(types.ConfigFlag, viper.ConfigFileUsed())
}
