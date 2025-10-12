package main

import (
	"log/slog"

	registry "github.com/Vahsek/distrokv/internal/registry"
	logging "github.com/Vahsek/distrokv/logging"
)

func main() {
	var fileLoggerProvider *logging.FileLoggerProvider = logging.NewFileLoggerProvider("logfile", ".log", "./logsdir", 10*1024, 10000)
	var logger slog.Logger = *logging.GetLogger(fileLoggerProvider)

	registry.StartRegistryServer(":8080", logger)
}
