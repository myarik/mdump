package storage

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type localStorage struct {
	filepath string
}

func (s localStorage) Save(ctx context.Context, fileName string, src io.Reader) error {
	dst, err := os.Create(filepath.Join(s.filepath, fileName))
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		if err = dst.Close(); err != nil {
			log.WithError(err).Warn("cannot close file")
		}
	}(dst)

	_, err = io.Copy(dst, src)
	return err
}

func NewLocalStorage(filepath string) *localStorage {
	return &localStorage{filepath}
}
