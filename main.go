package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var ignoreFiles = []string{".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/"}

type Contributions map[string]int

var contributions = make(Contributions)

func main() {
	process := make(chan string)
	authorChan := make(chan []string, 10)
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
		case authors := <-authorChan:
			countContribution(authors)
		}
	}

	printContributors(contributions.sort())
}

func scanFolder(root string, process chan string) {
	files, err := ioutil.ReadDir(root)

	if err != nil {
		log.Fatal(err.Error())
	}

LOOP:
	for _, file := range files {
		var filePath string
		if root == "./" {
			filePath = root + file.Name()
		} else {
			filePath = root + "/" + file.Name()
		}

		for _, ignore := range ignoreFiles {
			if strings.Contains(filePath, ignore) {
				continue LOOP
			}
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

func checkAuthor(filePath string, process chan []string) {
	if filePath == "DONE" {
		close(process)
		return
	}

	cmd := exec.Command("bash", "-c", "git blame "+filePath+" --porcelain | grep '^author ' | sort -u")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("file: ", filePath)
	process <- parseAuthor(string(output))
}

func parseAuthor(outputMessage string) []string {
	authors := strings.Split(strings.TrimSpace(outputMessage), "author ")
	fmt.Println(authors)
	return authors
}

func countContribution(authors []string) {
	for _, author := range authors {
		if author == "" {
			continue
		}
		trimmedAuthor := strings.TrimSpace(author)

		if _, ok := contributions[trimmedAuthor]; ok {
			contributions[trimmedAuthor]++
		} else {
			contributions[trimmedAuthor] = 1
		}
	}
}

//
//func sliceContain(searchValue string, slice []string) bool {
//	for _, item := range slice {
//		if item == searchValue {
//			return true
//		}
//	}
//	return false
//}

type sortStruct struct {
	Key   string
	Value int
}

func (array Contributions) sort() []sortStruct {

	var collection []sortStruct
	for key, value := range array {
		collection = append(collection, sortStruct{key, value})
	}

	sort.Slice(collection, func(i, j int) bool {
		return collection[i].Value > collection[j].Value
	})

	return collection
}

func printContributors(contributions []sortStruct) {
	for _, contributor := range contributions {
		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
	}
}
