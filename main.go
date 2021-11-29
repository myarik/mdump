package main

import (
	"github.com/myarik/mdump/cmd"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

var AppVersion = "unknown"

func main() {
	cmd.Execute()
}
