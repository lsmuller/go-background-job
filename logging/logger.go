package logging

import (
	"github.com/lsmuller/go-background-job/config"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/logger"
)

func CreateLogger(config *config.LoggingConfig) logrus.FieldLogger {
	log := logrus.New()
	log.SetReportCaller(true)

	// enable report called (file and line)
	switch config.Verbose {
	case 0:
		log.Level = logrus.ErrorLevel
	case 1:
		log.Level = logrus.WarnLevel
	case 2:
		log.Level = logrus.InfoLevel
	case 3:
		log.Level = logrus.DebugLevel
	default:
		log.Level = logrus.InfoLevel
	}

	if config.LogJSON {
		log.Formatter = new(logrus.JSONFormatter)
	}

	logger.SetLogger(log)
	return log
}
