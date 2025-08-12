package ulog

import (
	"os"
	"path/filepath"
)

var defaultLogFile = func() string {
	logPath := filepath.Join("logs", "app.log")
	err := os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		panic(err)
	}

	return logPath
}()
