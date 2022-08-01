package main

import (
	"fmt"
	"math"
	"regexp"
	"time"
)

type Contributions map[string]int
type LineResultChannel chan

//var contributions = make(Contributions)

var result []string

func main() {
	start := time.Now()
	filesChan := make(chan string)
	var lineChans chan (<-chan string)
	lineChans = make(chan (<-chan string), 50)

	go scanFolder("./", filesChan)

	go func() {
		defer close(lineChans)
		for path := range filesChan {
			lineChans <- checkAuthor(path)
		}
	}()

	for lineChan := range lineChans {
		select {
		case line := <-lineChan:
			result = append(result, line)
		}
	}

	for lineChan := range lineChans {
		select {
		case line := <-lineChan:
			result = append(result, line)
		default:
		}
	}

	fmt.Println("Total typescript files: ", len(result))
	var documentedFiles = filter(result, "NO DESCRIPTION FOUND")
	fmt.Println("Number of documentation files: ", len(documentedFiles))
	fmt.Println("Documentation coverage: ", math.Round(float64(len(documentedFiles))/float64(len(result))*100), " %")
	fmt.Printf("Executed in %s", time.Since(start))
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

func appendValidLine(path string) {

}

var descriptionRegex = "\\/\\*\\*"

//func printContributors(contributions []sortStruct) {
//	for _, contributor := range contributions {
//		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
//	}
//}
