package logging

import "sync"

type FileLoggerProvider struct {
	fileMetaData *FileMetaData
}

func NewFileLoggerProvider(filenamePrefix string, fileExtension string, fileDirectory string, maxFileSize int64, maxFileCount int32) *FileLoggerProvider {
	return &FileLoggerProvider{
		fileMetaData: &FileMetaData{
			filenamePrefix: filenamePrefix,
			fileExtension:  fileExtension,
			fileDirectory:  fileDirectory,
			maxFileSize:    maxFileSize,
			maxFileCount:   maxFileCount,
			fileLock:       sync.Mutex{},
			file:           nil,
		},
	}
}

func (fileLoggerProvider *FileLoggerProvider) Write(p []byte) (int, error) {
	var message string = string(p)
	fileLoggerProvider.fileMetaData.WriteLogsToFile(message)
	return len(message), nil
}
