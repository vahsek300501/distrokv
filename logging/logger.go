package logging

import (
	"io"
	"log/slog"
	"sync"
)

var (
	loggerInstance *slog.Logger
	once           sync.Once
)

func GetLogger(writers ...io.Writer) *slog.Logger {
	once.Do(func() {
		multiWriters := io.MultiWriter(writers...)
		loggerInstance = slog.New(slog.NewTextHandler(multiWriters, nil))
	})
	return loggerInstance
}
