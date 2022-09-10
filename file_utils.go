package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

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

func checkLines(filePath string) <-chan string {
	lineChan := make(chan string, 50)

	fileLines, err := readLines(filePath)

	if err != nil {
		defer close(lineChan)
		return lineChan
	}

	go func() {
		defer close(lineChan)
		regex := regexp.MustCompile(descriptionRegex)
		for _, line := range fileLines {
			descriptionLine := regex.FindString(line)
			if descriptionLine != "" {
				lineChan <- filePath
				return
			}
		}

		lineChan <- "NO DESCRIPTION FOUND"
	}()

	return lineChan
}

var descriptionRegex = "\\/\\*\\*"
