package main

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func Logger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)                  // default
	logger.Formatter.(*logrus.TextFormatter).DisableColors = true // remove colors
	logger.WithField("src", "ipsync")
	logger.SetLevel(logrus.DebugLevel)
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	if err == nil {
		logger.Hooks.Add(hook)
	}
	return logger
}

func log_err(msg string, fn string) {
	Logger().WithField("fn", fn).Error(msg)
}

func log_info(msg string, fn string) {
	Logger().WithField("fn", fn).Debug(msg)
}
