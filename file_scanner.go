package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LineParser struct {
	ignoreFiles []string
}

type LineProcessor struct {
	method string
}

//TODO: make this as args for many language type
func (parser *LineParser) construct() {
	parser.ignoreFiles = []string{".gitignore", ".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/", "allure", "artifacts"}
}

func (parser *LineParser) scanFolder(root string, fileChan chan<- string) {
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

		for _, ignore := range parser.ignoreFiles {
			if strings.Contains(filePath, ignore) {
				continue FILES_LOOP
			}
		}

		if file.IsDir() {
			parser.scanFolder(filePath, fileChan)
		} else {
			fileChan <- filePath
		}
	}

	if root == "./" {
		close(fileChan)
	}
}

func (lineProcessor *LineProcessor) execute(filesChan <-chan string) {
	if lineProcessor.method == "line" {
		countLines(filesChan)
	} else {
		countFiles(filesChan)
	}
}

func countFiles(filesChan <-chan string) {
	count := 0
	for range filesChan {
		count++
	}
	fmt.Println(count)
}

func countLines(fileChan <-chan string) {
	var result = 0
	lineChan := make(chan []string, 50)

	go func() {
		defer close(lineChan)
		for filePath := range fileChan {
			lineCount, err := readLines(filePath)
			if err != nil {
				continue
			} else {
				lineChan <- lineCount
			}
		}
	}()

	for lines := range lineChan {
		result = result + len(lines)
	}

	fmt.Printf("Total line of code: %d \n", result)
}
