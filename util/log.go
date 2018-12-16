package util

import (
	"os"

	logging "github.com/op/go-logging"
)

// InitLogger initializes the global logger
func InitLogger() logging.Logger {
	log := logging.MustGetLogger("log")
	logFormat := logging.MustStringFormatter(
		`%{color}%{time:[2006-01-02 15:04:05]} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logBackendFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	logging.SetBackend(logBackend, logBackendFormatter)

	return *log
}
