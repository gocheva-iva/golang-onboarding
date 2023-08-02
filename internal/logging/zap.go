package logging

import (
	"io/fs"
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	once   sync.Once
	Logger Log
)

type ZapLogger struct {
	logger *zap.Logger
}

func newZapLogger() Log {
	l, _ := zap.NewProduction()
	return &ZapLogger{logger: l}
}

func GetLogger() Log {
	once.Do(func() {
		Logger = newZapLogger()
	})
	return Logger
}

func (z *ZapLogger) FileLog(filePath string, msg string) {
	octalValue := 0644
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.FileMode(octalValue))
	if err != nil {
		GetLogger().Error("Error opening log file", err)
		return
	}
	defer f.Close()
	if _, err = f.Write([]byte(msg)); err != nil {
		GetLogger().Error("Error writing to log file", err)
		return
	}
}

func (z *ZapLogger) Info(msg string) {
	z.logger.Info(msg)
}

func (z *ZapLogger) Error(msg string, err ...error) {
	z.logger.Error(msg, zap.Error(err[0]))
}
