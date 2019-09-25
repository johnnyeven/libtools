package ww

import (
	"fmt"
)

type logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Notice(format string, v ...interface{})
}

var mainLogger logger = nil

func SetLogger(l logger) {
	mainLogger = l
}

func Debug(format string, v ...interface{}) {
	if mainLogger != nil {
		mainLogger.Debug(format, v...)
		return
	}
	fmt.Println(fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	if mainLogger != nil {
		mainLogger.Info(format, v...)
		return
	}
	fmt.Println(fmt.Sprintf(format, v...))
}

func Warning(format string, v ...interface{}) {
	if mainLogger != nil {
		mainLogger.Warning(format, v...)
		return
	}
	fmt.Println(fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	if mainLogger != nil {
		mainLogger.Error(format, v...)
		return
	}
	fmt.Println(fmt.Sprintf(format, v...))
}

func Notice(format string, v ...interface{}) {
	if mainLogger != nil {
		mainLogger.Notice(format, v...)
		return
	}
	fmt.Println(fmt.Sprintf(format, v...))
}
