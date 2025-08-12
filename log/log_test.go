package ulog

import (
	"github.com/labstack/gommon/log"
	"io"
	"os"
	"testing"
)

var logger = func() *log.Logger {
	l := log.New("pkg-tool")
	l.SetPrefix("pkg")
	lf, err := os.OpenFile("pkg_tool.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l.SetOutput(io.MultiWriter(
		lf,
		os.Stdout,
	))

	return l
}()

func TestLog(t *testing.T) {
	logger.Debug("debug info")
	logger.Info("info info")

	logger.Error("error info")

	logger.Warn("warn info")

}
