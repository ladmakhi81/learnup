package logger

type Log interface {
	Print(message LogMessage)
	Error(message LogMessage)
	Warning(message LogMessage)
}
