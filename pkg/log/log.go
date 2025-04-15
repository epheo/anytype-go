package log

// Level represents the severity level of a log message
type Level int

const (
	// LevelError only shows error messages
	LevelError Level = iota
	// LevelInfo shows errors and info messages
	LevelInfo
	// LevelDebug shows all messages including debug information
	LevelDebug
)

// String returns the string representation of a log level
func (l Level) String() string {
	switch l {
	case LevelError:
		return "error"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	default:
		return "unknown"
	}
}

// ParseLevel converts a string to a Level
func ParseLevel(level string) Level {
	switch level {
	case "error":
		return LevelError
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	default:
		return LevelInfo // default to info level
	}
}

// Logger defines the interface for logging operations
type Logger interface {
	Error(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	SetLevel(level Level)
	GetLevel() Level
}
