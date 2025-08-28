package logger

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)

func Log(level LogLevel, format string, a ...any) {
	prefix := ""
	switch level {
	case INFO:
		prefix = "[INFO] "
	case WARN:
		prefix = "[WARN] "
	case ERROR:
		prefix = "[ERROR] "
	}
	fmt.Printf(prefix+format+"\n", a...)
}

func Fatal(format string, a ...interface{}) {
    Log(ERROR, format, a...)
    os.Exit(1)
}
