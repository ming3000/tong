package common

import "io"

type Logger interface {
	SetOutput(w io.Writer)
	SetPrefix(p string)
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
