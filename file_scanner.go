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

type LineProcessor struct{}

func (parser *LineParser) construct() {
	parser.ignoreFiles = []string{".gitignore", ".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/"}
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
			continue FILES_LOOP
		}
	}

	if root == "./" {
		close(fileChan)
	}
}

func (lineProcessor *LineProcessor) execute(filesChan <-chan string) {
	//TODO: multiple method
	count := 0
	for range filesChan {
		count++
	}
	fmt.Println(count)
}

//TODO: count total line of code
//func countLines(path string) {
//
//}
//var descriptionRegex = "\\/\\*\\*"
//func printContributors(contributions []sortStruct) {
//	for _, contributor := range contributions {
//		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
//	}
//}
