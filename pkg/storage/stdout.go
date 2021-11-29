package storage

import (
	"io"
	"os"
)

type stdout struct{}

func (s stdout) Save(src io.Reader) error {
	_, err := io.Copy(os.Stdout, src)
	return err
}

func NewStdout() *stdout {
	return &stdout{}
}
