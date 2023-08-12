package logger

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/Code-Hex/dd"
	"github.com/kodmain/kitsune/src/internal/observability/logger/levels"
	"github.com/kodmain/kitsune/src/internal/observability/logger/writers"
)

const PATH = "| \033[38;5;%sm%s\033[39;49m |"

var instance *Logger = nil

type Logger struct {
	level   levels.TYPE
	success *log.Logger
	failure *log.Logger
}

func Default() *Logger {
	if instance == nil {
		instance = New(writers.DEFAULT, levels.DEFAULT)
	}
	return instance
}

func (l *Logger) Write(level levels.TYPE, messages ...any) {
	for _, message := range messages {
		var logger *log.Logger = nil
		if level <= levels.WARN {
			logger = l.failure
		} else {
			logger = l.success
		}

		if level == levels.DEBUG {
			logger.Println(fmt.Sprintf(PATH, level.Color(), level.String()), dd.Dump(message, dd.WithIndent(4)))
		} else if level <= l.level {
			logger.Println(fmt.Sprintf(PATH, level.Color(), level.String()), message)
		}
	}
}

func (l *Logger) Panic(err error) bool {
	if err != nil {
		l.Write(levels.PANIC, err, string(debug.Stack()))
		return true
	}

	return false
}

func (l *Logger) Fatal(err error) bool {
	if err != nil {
		l.Write(levels.FATAL, err, string(debug.Stack()))
		return true
	}

	return false
}

func (l *Logger) Error(err error) bool {
	if err != nil {
		l.Write(levels.ERROR, err)
		return true
	}

	return false
}

func (l *Logger) Success(v ...any) {
	l.Write(levels.SUCCESS, v...)
}

func (l *Logger) Message(v ...any) {
	l.Write(levels.MESSAGE, v...)
}

func (l *Logger) Warn(v ...any) {
	l.Write(levels.WARN, v...)
}

func (l *Logger) Info(v ...any) {
	l.Write(levels.INFO, v...)
}

func (l *Logger) Debug(v ...any) {
	l.Write(levels.DEBUG, v...)
}

func (l *Logger) Trace() {
	l.Write(levels.TRACE, string(debug.Stack()))
}

func New(t writers.TYPE, l levels.TYPE) *Logger {
	return &Logger{
		level:   l,
		success: log.New(writers.Make(t, writers.SUCCESS), "", log.Ldate|log.Ltime),
		failure: log.New(writers.Make(t, writers.FAILURE), "", log.Ldate|log.Ltime),
	}
}
