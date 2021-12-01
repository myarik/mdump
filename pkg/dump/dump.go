package dump

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, fileName string, src io.Reader) error
}

//NewPgDump returns a
func NewPgDump() *pgDumpService {
	return &pgDumpService{
		command: "pg_dump",
		options: []string{"-Fc", "--clean", "--no-acl", "--no-owner"},
	}
}
