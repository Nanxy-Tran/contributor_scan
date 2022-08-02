package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"time"
)

//type Contributions map[string]int
type LineResult struct {
	FilePath  string
	Author    string
	RawString string
}

//var contributions = make(Contributions)

type LineCrawler interface {
	scan(filePath string) <-chan LineResult
}

type GitCrawler struct {
	config interface{}
}

type RawStringCrawler struct {
	name      string
	predicate func()
}

func (scanner *GitCrawler) scan(filePath string) <-chan LineResult {
	authorChan := make(chan LineResult, 50)

	go func() {
		defer close(authorChan)
		cmd := exec.Command("bash", "-c", "git blame "+filePath+" --porcelain | sed 's/author //p' | sort | uniq -c| sort -rn | head -n 1 | sed 's/[0-9]*//g'")
		output, err := cmd.CombinedOutput()

		fmt.Println("Git blaming.... ", filePath)
		if err != nil {
			log.Fatal(err)
			return
		}
		authorChan <- parseAuthor(string(output), filePath)
	}()
	return authorChan
}

func (scanner *RawStringCrawler) scan(filePath string) <-chan LineResult {
	lineChan := make(chan LineResult, 50)

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
				lineChan <- LineResult{FilePath: filePath}
				return
			}
		}

		lineChan <- LineResult{FilePath: filePath}
	}()

	return lineChan
}

//TODO: Init the predicate
func LineCrawlerFactory(crawler string) LineCrawler {
	if crawler == "git" {
		return &GitCrawler{config: "git"}
	} else {
		return &RawStringCrawler{name: "line"}
	}
}

var result []LineResult

func main() {
	start := time.Now()
	filesChan := make(chan string)
	var lineChans chan (<-chan LineResult)
	lineChans = make(chan (<-chan LineResult), 50)

	go scanFolder("./", filesChan)

	go func() {
		defer close(lineChans)
		for path := range filesChan {
			fmt.Println("Checking.... ", path)
			lineChans <- checkLines(path)
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

func appendValidLine(path string) {

}

var descriptionRegex = "\\/\\*\\*"

//func printContributors(contributions []sortStruct) {
//	for _, contributor := range contributions {
//		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
//	}
//}
