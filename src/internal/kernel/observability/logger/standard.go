package logger

import (
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/writers"
)

var instance *Logger = nil

func standard() *Logger {
	if instance == nil {
		instance = New(writers.DEFAULT, levels.DEFAULT)
	}
	return instance
}

func SetLevel(l levels.TYPE) {
	standard().level = l
}

func Panic(err error) bool {
	return standard().Panic(err)
}

func Fatal(err error) bool {
	return standard().Fatal(err)
}

func Error(err error) bool {
	return standard().Error(err)
}

func Success(v ...any) {
	standard().Success(v...)
}

func Message(v ...any) {
	standard().Message(v...)
}

func Warn(v ...any) {
	standard().Warn(v...)
}

func Info(v ...any) {
	standard().Info(v...)
}

func Debug(v ...any) {
	standard().Debug(v...)
}

func Trace() {
	standard().Trace()
}
