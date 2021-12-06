package dump

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

const DATA10MB = int64(10 * 1024 * 1024)
const DATA1MB = int64(1024 * 1024)

func getDatabaseNameFromURI(dbURI string) string {
	return dbURI[strings.LastIndex(dbURI, "/")+1:]
}

// WriteCounter counts the number of bytes written to it.
type WriteCounter struct {
	Total        int64 // Total # of bytes transferred
	prevLogPoint int64
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	if wc.Total > (wc.prevLogPoint + DATA10MB) {
		wc.prevLogPoint = wc.Total
		log.Infof("%dMb uploaded ...", wc.Total/DATA1MB)
	}
	return n, nil
}

func (wc *WriteCounter) TotalSize() int64 {
	return wc.Total / DATA1MB
}
