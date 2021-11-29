package dump

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"sync"
	"time"
)

type dumpStorage interface {
	Save(dst io.Reader) error
}

type pgDumpService struct {
	command string
	options []string
}

func (s pgDumpService) Run(ctx context.Context, storage dumpStorage, credentials PGCredentials) error {
	if _, err := exec.LookPath(s.command); err != nil {
		log.Errorf("cannot find a command %s", s.command)
		return errors.New("cannot find the pg_dump command")
	}

	var cmdOptions []string
	cmdOptions = append(cmdOptions, s.options...)
	cmdOptions = append(cmdOptions, []string{"--dbname", credentials.String()}...)

	log.Debug("stating dump...")

	cmd := exec.CommandContext(ctx, s.command, s.options...)
	var errGzip, errStorag error

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "cannot create stdout pipe")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// read from stdout
	go func() {
		defer wg.Done()

		pr, pw := io.Pipe()
		gw := gzip.NewWriter(pw)
		gw.ModTime = time.Now()

		go func() {
			_, errGzip = io.Copy(gw, bufio.NewReader(outPipe))
			// close gzip
			if err = gw.Close(); err != nil {
				log.WithError(err).Warn("cannot close gzip writer")
			}
			if err = pw.Close(); err != nil {
				log.WithError(err).Warn("cannot close pipe writer")
			}
		}()
		errStorag = storage.Save(pr)
	}()

	// capture and read stderr in case an error occurs
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "cannot start a command")
	}
	// wait reading
	wg.Wait()

	// cmd.Wait() should be called only after we finish reading from outPipe
	// wg ensures that we finish
	if err = cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("dump canceled")
		}
		return errors.Wrap(err, "dump command returns an error")
	}
	if errStorag != nil {
		return errors.Wrap(err, "storage returns an error")
	}
	if errGzip != nil {
		return errors.Wrap(err, "gzip cannot read from stdout")
	}
	return nil
}

func New() *pgDumpService {
	return &pgDumpService{
		command: "pg_dump",
		options: []string{"-Fc", "--clean", "--no-acl", "--no-owner"},
	}
}
