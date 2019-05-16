package cmd

import (
	"fmt"
	"github.com/affiliator/mgw/config"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mgw",
	Short: "A simple smtp server which redirects every email to Mailgun.",
	Run:   nil,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		err := config.Ptr().Reload()
		if err != nil {
			panic(err)
		}
	})

	cfg := config.Ptr()

	rootCmd.PersistentFlags().StringVarP(&cfg.Paths.Config.Name, "config", "c",
		"", "Path to the configuration file")

	rootCmd.PersistentFlags().StringVarP(&cfg.Paths.Pid.Name, "pid", "p",
		"", "Path to the pid file")

	rootCmd.PersistentFlags().StringVarP(&cfg.Paths.Credentials.Name, "credentials", "a",
		"", "Path to the file containing mailgun api credentials")
}
