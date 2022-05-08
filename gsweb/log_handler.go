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

	logFileName := now.Format("2006-01-02") + ".log"
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
		TimestampFormat: "2006/01/02 - 15:04:05",
	})
	return logger
}

type LogField logrus.Fields

func LogInfo(msg string, fields *LogField) {
	if fields != nil {
		logrus.Info(msg, fields)                              // output in terminal
		Logger().WithFields(logrus.Fields(*fields)).Info(msg) // output in file
		return
	}
	logrus.Info(msg)
	Logger().Info(msg)
}

func LogError(msg string, fields *LogField) {
	if fields != nil {
		logrus.Error(msg, fields)
		Logger().WithFields(logrus.Fields(*fields)).Error(msg)
		return
	}
	logrus.Error(msg)
	Logger().Error(msg)
}

func LogWarn(msg string, fields *LogField) {
	if fields != nil {
		logrus.Warn(msg, fields)
		Logger().WithFields(logrus.Fields(*fields)).Warn(msg)
		return
	}
	logrus.Warn(msg)
	Logger().Warn(msg)
}

func LogDebug(msg string, fields *LogField) {
	if fields != nil {
		logrus.Debug(msg, fields)
		Logger().WithFields(logrus.Fields(*fields)).Debug(msg)
		return
	}
	logrus.Debug(msg)
	Logger().Debug(msg)
}
