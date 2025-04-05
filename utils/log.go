package utils

import (
	"fmt"
	"os"
	"path"
)

func SaveMessageIntoLog(logFileName string, message string) error {
	rootDir, rootDirErr := os.Getwd()
	if rootDirErr != nil {
		return rootDirErr
	}
	logFileDirectoryPath := path.Join(rootDir, "log")
	logFilePath := path.Join(logFileDirectoryPath, logFileName)

	if err := os.MkdirAll(logFileDirectoryPath, os.ModePerm); err != nil {
		return err
	}

	logFile, logFileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if logFileErr != nil {
		return logFileErr
	}
	defer logFile.Close()
	if _, err := logFile.WriteString(fmt.Sprintf("%s\n", message)); err != nil {
		return err
	}
	return nil
}
