package main

import (
	"log"
	"os"
	"strings"
)

var ignoreFiles = []string{".gitignore", ".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/"}

//var typescriptFileExtensions = []string{".tsx", ".ts"}

type FileCrawler struct {
	ignoreFiles        []string
	specificExtensions []string
	fileChannels       <-chan string
}

func (scanner *FileCrawler) scan(root string) {
	fileChannels := make(chan string)
	go scanner.scanFolder(root, fileChannels)

	scanner.fileChannels = fileChannels
}

func (scanner *FileCrawler) scanFolder(root string, fileChan chan<- string) {
	files, err := os.ReadDir(root)

	if err != nil {
		log.Fatal(err.Error())
	}

FILES_LOOP:
	for _, file := range files {
		var filePath = root + "/" + file.Name()
		for _, ignore := range scanner.ignoreFiles {
			if strings.Contains(filePath, ignore) {
				continue FILES_LOOP
			}
		}

		if file.IsDir() {
			scanner.scanFolder(filePath, fileChan)
		} else if scanner.specificExtensions != nil {
			for _, ext := range scanner.specificExtensions {
				if strings.Contains(filePath, ext) {
					fileChan <- filePath
					continue FILES_LOOP
				}
			}
		} else {
			fileChan <- filePath
			continue FILES_LOOP
		}
	}

	if root == "." {
		close(fileChan)
	}
}
