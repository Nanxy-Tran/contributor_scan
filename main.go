package main

import (
	"fmt"
	"regexp"
	"time"
)

type Contributions map[string]int
type LineResultChannel chan LineResult

//var contributions = make(Contributions)
/**
Git owner
Git author contributions
Line counts
Graph ?
*/

func main() {
	start := time.Now()
	filesChan := make(chan string)
	parser := parserFactory("git")
	parser.construct()

	processor := fileProcessorFactory("owner")
	go parser.scanFolder("./", filesChan)

	processor.execute(filesChan)
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
