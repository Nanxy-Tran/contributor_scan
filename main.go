package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ignoreFiles = []string{".gitignore", ".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/"}

type Contributions map[string]int

var contributions = make(Contributions)

func main() {
	start := time.Now()
	filesChan := make(chan string)
	var authorsChan chan (<-chan []string)
	authorsChan = make(chan (<-chan []string), 50)

	go scanFolder("./", filesChan)

	go func() {
		defer close(authorsChan)
		for path := range filesChan {
			fmt.Println("Checking at path", path)
			authorsChan <- checkAuthor(path)
		}
	}()

	for authors := range authorsChan {
		select {
		case authorsString := <-authors:
			countContribution(authorsString)
		}
	}

	printContributors(contributions.sort())
	fmt.Printf("Executed in %s", time.Since(start))
}

func scanFolder(root string, fileChan chan<- string) {
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
			fileChan <- filePath
		}
	}

	if root == "./" {
		close(fileChan)
	}
}

func checkAuthor(filePath string) <-chan []string {
	authorChan := make(chan []string)
	go func() {
		defer close(authorChan)
		cmd := exec.Command("bash", "-c", "git blame "+filePath+" --porcelain | grep '^author ' | sort -u")
		output, err := cmd.CombinedOutput()
		fmt.Println("Git blame for ", filePath)
		if err != nil {
			log.Fatal(err)
			return
		}
		authorChan <- parseAuthor(string(output))
	}()
	return authorChan
}

func parseAuthor(outputMessage string) []string {
	authors := strings.Split(strings.TrimSpace(outputMessage), "author ")
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
