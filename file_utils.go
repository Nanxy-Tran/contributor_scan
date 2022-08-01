package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var ignoreFiles = []string{".gitignore", ".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/"}
var typescriptFileExtensions = []string{".tsx", ".ts"}

func scanFolder(root string, fileChan chan<- LineResultChannel) {
	files, err := ioutil.ReadDir(root)

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

		for _, ignore := range ignoreFiles {
			if strings.Contains(filePath, ignore) {
				continue FILES_LOOP
			}
		}

		if file.IsDir() {
			scanFolder(filePath, fileChan)
		} else {
			for _, ext := range typescriptFileExtensions {
				if strings.Contains(filePath, ext) {
					fileChan <- filePath
					continue FILES_LOOP
				}
			}
		}
	}

	if root == "./" {
		close(fileChan)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err, "error happened")
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
