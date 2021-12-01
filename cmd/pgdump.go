package cmd

import (
	"context"
	"github.com/myarik/mdump/pkg/dump"
	"github.com/myarik/mdump/pkg/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	cmdPgDump.PersistentFlags().String("pg_uri", "",
		"PostgreSQL connection URI (postgresql://user:secret@localhost/db_name)")
	if err := viper.BindPFlag(
		"pg_uri", cmdPgDump.PersistentFlags().Lookup("pg_uri")); err != nil {
		log.WithError(err).Panic("viper library returns an error")
	}

	cmdPgDumpLocal.Flags().String("path", "", "path to the dump location")
	if err := viper.BindPFlag(
		"local_path", cmdPgDumpLocal.Flags().Lookup("path")); err != nil {
		log.WithError(err).Panic("viper library returns an error")
	}

	cmdPgDumpS3.Flags().String("aws-bucket", "", "AWS bucket")
	if err := viper.BindPFlag(
		"aws_bucket", cmdPgDumpS3.Flags().Lookup("aws-bucket")); err != nil {
		log.WithError(err).Panic("viper library returns an error")
	}

	cmdPgDumpS3.Flags().String("aws-bucket-key", "", "AWS bucket's key")
	if err := viper.BindPFlag(
		"aws_bucket_key", cmdPgDumpS3.Flags().Lookup("aws-bucket-key")); err != nil {
		log.WithError(err).Panic("viper library returns an error")
	}
}

var cmdPgDumpS3 = &cobra.Command{
	Use:   "s3",
	Short: "a dump is stored in an Amazon S3 bucket",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("aws_bucket") == "" {
			return errors.New("required flag(s) \"aws-bucket\" not set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDump(
			storage.NewS3Storage(
				viper.GetString("aws_bucket"), viper.GetString("aws_bucket_key")),
			viper.GetString("pg_uri"),
		); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

var cmdPgDumpLocal = &cobra.Command{
	Use:   "local",
	Short: "a dump is stored locally",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("local_path") == "" {
			return errors.New("required flag(s) \"path\" not set")
		}
		if viper.GetString("pg_uri") == "" {
			return errors.New("required flag(s) \"pg_uri\" not set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDump(
			storage.NewLocalStorage(viper.GetString("local_path")),
			viper.GetString("pg_uri"),
		); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
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

func runDump(storage dump.Storage, dbURI string) error {
	s := dump.NewPgDump()
	return s.Run(context.Background(), storage, dbURI)
}
