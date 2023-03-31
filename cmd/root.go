package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/cmd/project"
	"os"
)

var (
	cfgFile   string
	baseUrl   string
	userToken string
	limit     int

	rootCmd = &cobra.Command{
		Use:   "gobit",
		Short: "gobit - enkel CLI for Bitbucket",
		Long:  `gobit lalala`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.gobit.yaml")
	rootCmd.PersistentFlags().StringVar(&baseUrl, "baseUrl", "http://git.spk.no", "base url for BitBucket instance")
	rootCmd.PersistentFlags().IntVar(&limit, "limit", 100, "max return values")
	rootCmd.PersistentFlags().StringVar(&userToken, "token", "", "token for user")

	viper.BindPFlag("baseUrl", rootCmd.PersistentFlags().Lookup("baseUrl"))
	viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))

	viper.SetDefault("baseUrl", "http://git.spk.no")
	viper.SetDefault("limit", 100)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(project.RootCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gobit")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
