package storage

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type fileStorage struct {
	filepath string
}

func (s fileStorage) Save(src io.Reader) error {
	dst, err := os.Create(s.filepath)
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

func NewFileStorage(filepath string) *fileStorage {
	return &fileStorage{filepath}
}
