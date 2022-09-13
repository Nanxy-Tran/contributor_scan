package main

import (
	"bufio"
	"fmt"
	"os"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err, "Can not read file: ", path)
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
	fmt.Println(len(fileLines))

	if err != nil {
		defer close(lineChan)
		return lineChan
	}

	go func() {
		defer close(lineChan)
		for _, line := range fileLines {
			lineChan <- line
		}
	}()

	return lineChan
}

//var descriptionRegex = "\\/\\*\\*"
