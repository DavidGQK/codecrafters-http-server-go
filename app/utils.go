package main

import (
	"os"
	"path/filepath"
)

func findFile(dir, filename string) (string, bool) {
	fPath := filepath.Join(dir, filename)
	content, err := os.ReadFile(fPath)
	if err != nil {
		return "", false
	}

	return string(content[:]), true
}

func saveFile(dir, filename string, fileContent []byte) error {
	fPath := filepath.Join(dir, filename)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(fileContent)
	if err != nil {
		return err
	}

	return nil
}
