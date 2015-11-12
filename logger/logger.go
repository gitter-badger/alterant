package logWrapper

import "log"

// LogWrapper formats log output specific to alterant
type LogWrapper struct {
	Verbose bool
}

// Info displays messages in green
func (l *LogWrapper) Info(format string, v ...interface{}) {
	format = "\033[32m\033[1m" + format + "\033[0m"
	if l.Verbose {
		log.Printf(format, v...)
	}
}

// NewLogWrapper returns an instance of `LogWrapper`
func NewLogWrapper(verbose bool) *LogWrapper {
	return &LogWrapper{Verbose: verbose}
}
