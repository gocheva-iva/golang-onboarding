package logging

type Log interface {
	FileLog(filePath string, msg string)
	Info(msg string)
	Error(msg string, err ...error)
}
