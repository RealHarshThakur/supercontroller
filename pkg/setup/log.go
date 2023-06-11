package setup

import (
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Logging sets up the logging for the router daemon
func Logging() *logrus.Logger {
	// Logging create logging object
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", fileName
		},
	})

	return log
}
