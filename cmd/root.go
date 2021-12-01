package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var mainCmd = &cobra.Command{
	Use:   "mdump",
	Short: "mdump is simple backup utility",
	Long:  "mdump allows creating a database dump",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				return
			}
		}
		os.Exit(0)
	},
}

func Execute() {
	viper.SetEnvPrefix("mdump")
	viper.AutomaticEnv()

	mainCmd.AddCommand(cmdPgDump)
	cmdPgDump.AddCommand(cmdPgDumpLocal)
	cmdPgDump.AddCommand(cmdPgDumpS3)

	if err := mainCmd.Execute(); err != nil {
		log.WithError(err).Error("unexpected error")
		os.Exit(1)
	}
}
