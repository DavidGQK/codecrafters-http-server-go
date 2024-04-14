package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func findFile(dir, filename string) (string, bool) {
	currPath, _ := os.Getwd()
	fmt.Println("curr_path:", currPath)

	fPath := filepath.Join(dir, filename)
	fmt.Println("filepath:", fPath)
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
