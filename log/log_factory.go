package ulog

import (
	"github.com/labstack/gommon/log"
	"io"
	"os"
)

func init() {
	log.EnableColor()
	log.SetLevel(log.DEBUG)
}

func NewSimpleLogger(loggerName string) (*log.Logger, error) {
	lf, err := os.OpenFile(defaultLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	l := log.New(loggerName)

	l.SetOutput(io.MultiWriter(
		lf,
		os.Stdout,
	))

	l.EnableColor()
	l.SetLevel(log.DEBUG)

	return l, nil
}
