package common

import (
	"testing"
	"time"
)

func TestNewDefaultLogger(t *testing.T) {
	l := NewDefaultLogger()
	l.Debug("haha", "xixi", "lala")
	l.DebugFormat("year:%d-moth:%d-day:%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
}

func firstlevel() {
	secondLevel()
}

func secondLevel() {
	thirdLevel()
}

func thirdLevel() {
	l := NewDefaultLogger()

	ll := l.WithCaller(2)
	l.Debug(ll.callers)

	lll := l.WithCallersFrames()
	l.Debug(lll.callers)
}

func TestLogger_Debug(t *testing.T) {
	firstlevel()
}
