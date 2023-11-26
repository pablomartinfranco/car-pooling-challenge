package logger

import (
	"log"
	"os"
)

func NewLogger(name string) *log.Logger {
	if name != "" {
		name = "[" + name + "] "
	}
	return log.New(os.Stdout, name, log.Ldate|log.Ltime)
}
