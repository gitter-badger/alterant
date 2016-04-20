package logWrapper

import (
	"log"
	"strings"
)

// LogWrapper formats log output specific to alterant
type LogWrapper struct {
	Verbose bool
}

// Info displays messages in green
func (l *LogWrapper) Info(level int, format string, v ...interface{}) {
	format = "\033[36m\033[1m" + "<" + strings.Repeat("=", level) + "> " + format + "\033[0m"
	if l.Verbose {
		log.Printf(format, v...)
	}
}

// NewLogWrapper returns an instance of `LogWrapper`
func NewLogWrapper(verbose bool) *LogWrapper {
	return &LogWrapper{Verbose: verbose}
}
