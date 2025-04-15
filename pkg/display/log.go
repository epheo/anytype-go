package display

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// LogLevelError only shows error messages
	LogLevelError LogLevel = iota
	// LogLevelInfo shows errors and info messages
	LogLevelInfo
	// LogLevelDebug shows all messages including debug information
	LogLevelDebug
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "error"
	case LogLevelInfo:
		return "info"
	case LogLevelDebug:
		return "debug"
	default:
		return "unknown"
	}
}

// ParseLogLevel converts a string to a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "error":
		return LogLevelError
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	default:
		return LogLevelInfo // default to info level
	}
}
