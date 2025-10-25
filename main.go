package main

import (
	"log/slog"
	"os"

	logging "github.com/Vahsek/distrokv/internal/logging"
	registry "github.com/Vahsek/distrokv/internal/registry"
	node_service "github.com/Vahsek/distrokv/internal/worker_node/service"
)

func main() {
	var fileLoggerProvider *logging.FileLoggerProvider = logging.NewFileLoggerProvider(
		"logfile",
		".log",
		"C:\\Users\\kgambhir\\OneDrive - Microsoft\\Desktop\\distrokv\\nodelogs",
		10*1024,
		10000)
	var logger slog.Logger = *logging.GetLogger(fileLoggerProvider, os.Stdout)
	var bootType int = 1

	if bootType == 0 {
		registry.StartRegistryServer(":8080", logger)
	} else {
		workerNodeService := node_service.InitializeNewNodeService(
			"localhost",
			"127.0.0.1",
			"8002",
			"9002",
			1,
			"127.0.0.1:8080",
			logger)
		workerNodeService.BootstrapWorkerNode()
	}

}
