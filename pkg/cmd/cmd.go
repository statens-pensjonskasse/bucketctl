package cmd

import (
	"bucketctl/pkg/cmd/apply"
	"bucketctl/pkg/cmd/config"
	"bucketctl/pkg/cmd/get"
	"bucketctl/pkg/cmd/git"
	"bucketctl/pkg/cmd/pullRequest"
	"bucketctl/pkg/cmd/version"
	"bucketctl/pkg/common"
	"bucketctl/pkg/printer"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	cfgFile      string
	context      string
	outputFormat printer.OutputFormatType

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, common.ConfigFlag, "", "Base config file (default $HOME/.bucketctl/config.yaml")
	rootCmd.PersistentFlags().StringVarP(&context, common.ContextFlag, common.ContextFlagShorthand, "", "Context to use for overriding base config")
	rootCmd.PersistentFlags().String(common.BaseUrlFlag, "", "Base url for BitBucket instance")
	rootCmd.PersistentFlags().String(common.GitUrlFlag, "", "Base url for Git-commands")
	rootCmd.PersistentFlags().IntP(common.LimitFlag, common.LimitFlagShorthand, 1000, "Max return values")
	rootCmd.PersistentFlags().StringP(common.TokenFlag, common.TokenFlagShorthand, "", "Http access token")
	rootCmd.PersistentFlags().VarP(&outputFormat, common.OutputFlag, common.OutputFlagShorthand, "Output format. One of: pretty, yaml, json")

	viper.BindPFlag(common.BaseUrlFlag, rootCmd.PersistentFlags().Lookup(common.BaseUrlFlag))
	viper.BindPFlag(common.GitUrlFlag, rootCmd.PersistentFlags().Lookup(common.GitUrlFlag))
	viper.BindPFlag(common.LimitFlag, rootCmd.PersistentFlags().Lookup(common.LimitFlag))
	viper.BindPFlag(common.TokenFlag, rootCmd.PersistentFlags().Lookup(common.TokenFlag))
	viper.BindPFlag(common.OutputFlag, rootCmd.PersistentFlags().Lookup(common.OutputFlag))

	rootCmd.AddCommand(apply.Cmd)
	rootCmd.AddCommand(get.Cmd)

	rootCmd.AddCommand(git.Cmd)
	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(pullRequest.Cmd)
	rootCmd.AddCommand(version.Cmd)
}

func initConfig() {
	cfgPath, err := common.GetConfigPath()
	cobra.CheckErr(err)
	viper.AddConfigPath(cfgPath)

	if cfgFile == "" {
		cfgFile = filepath.Join(cfgPath, "config.yaml")
	}

	cobra.CheckErr(common.CreateDirIfNotExists(cfgFile, 0700))
	cobra.CheckErr(common.CreateFileIfNotExists(cfgFile, 0600))
	cobra.CheckErr(common.CheckFilePermission(cfgFile, 0600))

	viper.AutomaticEnv()

	viper.SetConfigFile(cfgFile)
	viper.SetConfigPermissions(0600)
	cobra.CheckErr(viper.ReadInConfig())

	if context != "" {
		viper.SetConfigName(context)
		viper.SetConfigPermissions(0600)
		cobra.CheckErr(viper.MergeInConfig())
		viper.SetConfigFile(cfgFile)
	}
}
