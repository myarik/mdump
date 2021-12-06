package dump

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os/exec"
	"sync"
	"time"
)

type pgDumpService struct {
	command string
	options []string
}

func (s pgDumpService) Run(ctx context.Context, storage Storage, dbURI string) error {

	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(ctx, 2*time.Hour)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	if _, err := exec.LookPath(s.command); err != nil {
		log.Errorf("cannot find a command %s", s.command)
		return errors.New("cannot find the pg_dump command")
	}

	log.Info("creating a dump...")

	var cmdOptions []string
	cmdOptions = append(cmdOptions, s.options...)
	cmdOptions = append(cmdOptions, []string{"--dbname", dbURI}...)

	dbName := getDatabaseNameFromURI(dbURI)
	dumpName := fmt.Sprintf("%s_%s.dump.gz", dbName, time.Now().Format("2006_01_02__15_04"))

	cmd := exec.CommandContext(ctx, s.command, cmdOptions...)

	log.WithField("cmd", cmd.String()).Debug("dump command")

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "cannot create stdout pipe")
	}

	// capture and read stderr in case an error occurs
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start process
	if err = cmd.Start(); err != nil {
		return errors.Wrap(err, "cannot start a command")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	progressCounter := WriteCounter{}
	// read from stdout
	go func() {
		defer wg.Done()

		pr, pw := io.Pipe()
		gw := gzip.NewWriter(pw)
		gw.ModTime = time.Now()

		go func() {
			_, err = io.Copy(gw, bufio.NewReader(outPipe))
			if err != nil {
				log.WithError(err).Error("cannot gzip data")
				cancel()
			}
			// close gzip
			if err = gw.Close(); err != nil {
				log.WithError(err).Warn("cannot close gzip writer")
			}
			if err = pw.Close(); err != nil {
				log.WithError(err).Warn("cannot close pipe writer")
			}
		}()
		// Wrap it with our custom io.Reader.
		test := io.TeeReader(pr, &progressCounter)
		// save a dump
		if err = storage.Save(ctx, dumpName, test); err != nil {
			log.WithError(err).Error("storage returns an error")
			cancel()
		}
	}()

	// Read stderr
	stderr, errStderr := ioutil.ReadAll(errPipe)
	// wait reading
	wg.Wait()

	// cmd.Wait() should be called only after we finish reading from outPipe
	// wg ensures that we finish
	if err = cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.WithField("cmd", cmd.String()).Error("execution timeout")
			return errors.New("timeout")
		}
		if ctx.Err() == context.Canceled {
			log.WithError(err).WithField("cmd", cmd.String()).Error("dump canceled")
			return errors.New("dump canceled")
		}
		log.WithError(err).WithField("cmd", cmd.String()).Error("dump command returns an error")
		return errors.Wrap(err, "dump command returns an error")
	}

	if errStderr != nil {
		log.WithError(errStderr).WithField("database", dbName).Error("cannot read stderr")
	}
	if len(stderr) > 0 {
		log.WithField("pg_dump error", stderr).Error("pg_dump returns an error")
	}

	log.WithFields(log.Fields{
		"database": dbName,
		"dump":     dumpName,
		"size":     progressCounter.TotalSize(),
	}).Info("Dump created")
	return nil
}
