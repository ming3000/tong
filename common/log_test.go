package common

import (
	"testing"
	"time"
)

func firstlevel(l *Logger) {
	secondLevel(l)
}

func secondLevel(l *Logger) {
	thirdLevel(l)
}

func thirdLevel(l *Logger) {
	l.Error("haha", "xixi", "lala")
	l.ErrorFormat("year:%d-moth:%d-day:%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
}

func TestLogger_Debug(t *testing.T) {
	l := NewDefaultLogger()
	firstlevel(l)
}
