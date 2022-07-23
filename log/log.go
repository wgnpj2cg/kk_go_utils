package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func NewLogger(LogFilePath string) *logrus.Logger {
	var log = logrus.New()
	writer1 := os.Stdout

	writer2, err := os.OpenFile(LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("Failed to log to file, using default stderr")
	}
	log.SetOutput(io.MultiWriter(writer1, writer2))
	//log.SetReportCaller(true)
	log.Formatter = &logrus.JSONFormatter{}
	//log.SetLevel(logrus.DebugLevel)
	return log
}
