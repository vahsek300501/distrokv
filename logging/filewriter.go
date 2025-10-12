package logging

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FileMetaData struct {
	filenamePrefix string
	fileExtension  string
	fileDirectory  string
	maxFileSize    int64
	maxFileCount   int32
	fileLock       sync.Mutex
	file           *os.File
}

func (fileMetaData *FileMetaData) InitializeFileMetaData(filenamePrefix string, fileExtension string, fileDirectory string, maxFileSize int64, maxFileCount int32) {
	fileMetaData.filenamePrefix = filenamePrefix
	fileMetaData.fileExtension = fileExtension
	fileMetaData.fileDirectory = fileDirectory
	fileMetaData.maxFileSize = maxFileSize
	fileMetaData.maxFileCount = maxFileCount
	fileMetaData.file = nil
}

func (fileMetaData *FileMetaData) WriteLogsToFile(logMessage string) {
	fileMetaData.fileLock.Lock()
	defer fileMetaData.fileLock.Unlock()

	if fileMetaData.file == nil {
		fileMetaData.createNewFile()
	}

	fileStats, _ := fileMetaData.file.Stat()
	if fileStats.Size()+int64(len(logMessage)) >= fileMetaData.maxFileSize {
		fileMetaData.rotateFile()
	}

	_, err := fileMetaData.file.WriteString(logMessage + "\n")
	if err != nil {
		return
	}
}

func (fileMetaData *FileMetaData) rotateFile() {
	fileMetaData.createNewFile()
	files, _ := os.ReadDir(fileMetaData.fileDirectory)
	if len(files) >= int(fileMetaData.maxFileCount) {
		fileMetaData.evictOlderFile()
	}
}

func (fileMetaData *FileMetaData) createNewFile() {
	var cntTimestamp string = strconv.FormatInt(time.Now().Unix(), 10)
	var newFileName string = fileMetaData.filenamePrefix + cntTimestamp + fileMetaData.fileExtension
	var fileLocation string = filepath.Join(fileMetaData.fileDirectory, newFileName)

	// Creating a new file
	newFile, err := os.Create(fileLocation)
	if err != nil {
		panic(err)
	}

	// If there is already a file closing it
	if fileMetaData.file != nil {
		error := fileMetaData.file.Close()
		if error != nil {
			panic(error)
		}
	}

	// Updating the file pointer in the struct to the new file
	fileMetaData.file = newFile
}

func (fileMetaData *FileMetaData) evictOlderFile() {
	files, err := os.ReadDir(fileMetaData.fileDirectory)
	if err != nil {
		panic(err)
	}

	var logFiles []string
	for _, file := range files {
		fileNameSplit := strings.Split(file.Name(), ".")
		if "."+fileNameSplit[len(fileNameSplit)-1] == fileMetaData.fileExtension {
			logFiles = append(logFiles, file.Name())
		}
	}
	sort.Strings(logFiles)
	fmt.Println(logFiles)
	var evictedFilePath string = path.Join(fileMetaData.fileDirectory, logFiles[0])
	fmt.Println(evictedFilePath)
	os.Remove(evictedFilePath)
}
