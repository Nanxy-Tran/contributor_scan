package main

import (
	"log"
	"os"
	"strings"
)

const LINE_COUNTER = "LINE_COUNTER"

type FileCrawler struct {
	ignoreFiles        []string
	specificExtensions []string
	fileChannels       <-chan string
}

func (scanner *FileCrawler) scan(root string) {
	fileChannels := make(chan string, 10)
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
		var filePath string

		if root == "./" {
			filePath = root + file.Name()
		} else {
			filePath = root + "/" + file.Name()
		}

		for _, ignore := range scanner.ignoreFiles {
			if strings.Contains(filePath, ignore) {
				continue FILES_LOOP
			}
		}

		if file.IsDir() {
			scanFolder(filePath, fileChan)
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

	if root == "./" {
		close(fileChan)
	}
}
