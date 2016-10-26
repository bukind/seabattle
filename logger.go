package seabattle

import (
	"log"
	"os"
)

var out *log.Logger

func init() {
	out = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func SetLogger(thelog *log.Logger) {
	out = thelog
}

func Logger() *log.Logger {
	return out
}
