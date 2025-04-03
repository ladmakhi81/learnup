package logger

type Log interface {
	Print(LogMessage)
	Error(LogMessage)
	Warning(LogMessage)
}
