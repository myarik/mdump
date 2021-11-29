package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var echoTimes int

var cmdPgDumpS3 = &cobra.Command{
	Use:   "s3",
	Short: "a dump is stored in an Amazon S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.Info("save to s3")
	},
}

var cmdPgDumpLocal = &cobra.Command{
	Use:   "local",
	Short: "a dump is stored locally",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.Infof("save to file %s", viper.GetString("local_path"))
	},
}

var cmdPgDump = &cobra.Command{
	Use:   "pgdump",
	Short: "command for backing up a PostgreSQL database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				return
			}
		}
		os.Exit(0)
	},
}
