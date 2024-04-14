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
