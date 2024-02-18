package logger

type Logger interface {
	Logf(format string, args ...interface{})
}
