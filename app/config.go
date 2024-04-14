package main

import "flag"

type config struct {
	directory string
	dirExists bool
}

func getConfig() *config {
	var appConfig config
	flag.StringVar(&appConfig.directory, "directory", "", "Directory")

	flag.Parse()

	if len(appConfig.directory) > 0 {
		appConfig.dirExists = true
	}

	return &appConfig
}
