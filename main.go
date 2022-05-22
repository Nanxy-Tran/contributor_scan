package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var ignoreFiles = []string{".git", ".idea"}
var contributions = make(map[string]int)

func main() {
	process := make(chan string)
	authorChan := make(chan string, 10)
	go scanFolder("./", process)

ProcessFile:
	for {
		select {
		case fileLocation := <-process:
			if fileLocation == "DONE" {
				break ProcessFile
			} else {
				go checkAuthor(fileLocation, authorChan)
			}
		}

		select {
		case author := <-authorChan:
			if author == "DONE" {
				close(authorChan)
			} else {
				countContribution(author)
			}
		}
	}

	printContributors(contributions)
}

func scanFolder(root string, process chan string) {
	files, err := ioutil.ReadDir(root)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		var filePath string
		if root == "./" {
			filePath = root + file.Name()
		} else {
			filePath = root + "/" + file.Name()
		}

		if sliceContain(file.Name(), ignoreFiles) {
			continue
		}

		if file.IsDir() {
			scanFolder(filePath, process)
			continue
		}

		process <- filePath
	}

	if root == "./" {
		process <- "DONE"
	}

	return
}

func checkAuthor(filePath string, process chan string) {
	if filePath == "DONE" {
		process <- "DONE"
		return
	}

	cmd := exec.Command("bash", "-c", "git blame "+filePath+" --porcelain | grep '^author ' | sort -u")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("file in checkAuthor", filePath)
	process <- parseAuthor(string(output))

}

func parseAuthor(outputMessage string) string {
	author := strings.Split(outputMessage, "author ")[1]
	return author
}

func countContribution(author string) {
	if _, ok := contributions[author]; ok {
		contributions[author]++
	} else {
		contributions[author] = 1
	}

}

func sliceContain(searchValue string, slice []string) bool {
	for _, item := range slice {
		if item == searchValue {
			return true
		}
	}
	return false
}

func printContributors(contributions map[string]int) {
	for contributor := range contributions {
		fmt.Println("Contributor: " + contributor + " has contributed " + strconv.Itoa(contributions[contributor]) + " files")
	}
}
