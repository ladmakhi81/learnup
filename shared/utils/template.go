package utils

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

func ParseTemplate(templatePath string, data any) (string, error) {
	rootDir, rootDirErr := os.Getwd()
	if rootDirErr != nil {
		return "", rootDirErr
	}
	path := filepath.Join(rootDir, "templates", templatePath)
	temp, tempErr := template.ParseFiles(path)
	if tempErr != nil {
		return "", tempErr
	}
	var tempBuffer bytes.Buffer
	if err := temp.Execute(&tempBuffer, data); err != nil {
		return "", err
	}
	return tempBuffer.String(), nil
}
