package common

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"runtime"
)

type Logger struct {
	stdLogger *log.Logger
	debug     bool
	callers   []string
}

func NewDefaultLogger() *Logger {
	return NewLogger("tong.log", 4, 1, "tong says:", true)
}

func NewLogger(fileName string, fileMaxSize int, fileMaxExpire int, prefix string, debug bool) *Logger {
	lj := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    fileMaxSize,
		MaxAge:     fileMaxExpire,
		MaxBackups: 1,
	}
	stdLogger := log.New(io.MultiWriter(os.Stdout, lj),
		prefix,
		log.Ldate|log.Ltime|log.Lshortfile)
	return &Logger{
		stdLogger: stdLogger,
		debug:     debug,
	}
}

func (l *Logger) copy() *Logger {
	cp := *l
	return &cp
}

func (l *Logger) WithCaller(skipLevel int) *Logger {
	cp := l.copy()
	pc, file, line, ok := runtime.Caller(skipLevel)
	if ok {
		f := runtime.FuncForPC(pc)
		callerInfo := fmt.Sprintf("%s::%d::%s", file, line, f.Name())
		cp.callers = []string{callerInfo}
	} // of>
	return cp
}

func (l *Logger) WithCallersFrames() *Logger {
	minDepth, maxDepth := 1, 7
	callers := make([]string, 0, maxDepth)
	pcs := make([]uintptr, maxDepth)

	depth := runtime.Callers(minDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callerInfo := fmt.Sprintf("%s::%d::%s", frame.File, frame.Line, frame.Function)
		callers = append(callers, callerInfo)
		if !more {
			break
		} //>>
	} // for>

	cp := l.copy()
	cp.callers = callers
	return cp
}

func (l *Logger) DebugFormat(format string, message ...interface{}) {
	if l.debug {
		l.stdLogger.Println(fmt.Sprintf(format, message...))
	} // if>
}

func (l *Logger) Debug(message ...interface{}) {
	if l.debug {
		l.stdLogger.Println(message...)
	} // if>
}

func (l *Logger) ErrorFormat(format string, message ...interface{}) {
	l.WithCallersFrames().stdLogger.Println(fmt.Sprintf(format, message...))
}

func (l *Logger) Error(format string, message ...interface{}) {
	l.WithCallersFrames().stdLogger.Println(message...)
}
