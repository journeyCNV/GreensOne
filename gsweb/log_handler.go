package gsweb

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func Logger() *logrus.Logger {
	now := time.Now()
	logFilePath := ""

	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}

	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}

	logFileName := now.Format("2022-01-01") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}

	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	logger := logrus.New()

	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

type LogField logrus.Fields

func LogInfo(msg string, fields *LogField) {
	logrus.Info(msg)
	if fields != nil {
		Logger().WithFields(logrus.Fields(*fields)).Info(msg)
		return
	}
	Logger().Info(msg)
}

func LogError(msg string, fields *LogField) {
	logrus.Error(msg)
	if fields != nil {
		Logger().WithFields(logrus.Fields(*fields)).Error(msg)
		return
	}
	Logger().Error(msg)
}

func LogWarn(msg string, fields *LogField) {
	logrus.Warn(msg)
	if fields != nil {
		Logger().WithFields(logrus.Fields(*fields)).Warn(msg)
		return
	}
	Logger().Warn(msg)
}

func LogDebug(msg string, fields *LogField) {
	logrus.Debug(msg)
	if fields != nil {
		Logger().WithFields(logrus.Fields(*fields)).Debug(msg)
		return
	}
	Logger().Debug(msg)
}
