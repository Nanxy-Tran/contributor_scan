package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type LineResult struct {
	FilePath  string
	Author    string ""
	RawString string
	Error     error
}

type GitParser struct {
	method string
}

func (parser *GitParser) execute(fileChanel <-chan string) {
	if parser.method == "owner" {
		lineChannel := make(chan (<-chan LineResult), 10)

		go func() {
			defer close(lineChannel)
			for path := range fileChanel {
				lineChannel <- checkAuthor(path)
			}
		}()

		var result []LineResult

		for lineChan := range lineChannel {
			select {
			case line := <-lineChan:
				result = append(result, line)
			}
		}

		for lineChan := range lineChannel {
			select {
			case line := <-lineChan:
				result = append(result, line)
			default:
			}
		}

		printGitResult(result)
		return
	}

	if parser.method == "contribute" {
		//TODO:
	}
}

func checkAuthor(filePath string) <-chan LineResult {
	authorChan := make(chan LineResult)

	go func() {
		defer close(authorChan)
		cmd := exec.Command("bash", "-c", "git blame "+filePath+" --line-porcelain | sed -n -e 's/author //p' | sort | uniq -c| sort -rn | head -n 2 | sed 's/[0-9]*//g'")
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Println("Git error: ", err.Error())
			return
		}
		authorChan <- parseAuthors(string(output), filePath)
	}()

	return authorChan
}

func parseAuthors(output string, filePath string) LineResult {
	var outputTrimmed = strings.TrimSpace(output)
	var outputArr = strings.Split(outputTrimmed, "  ")

	var result = LineResult{
		FilePath: filePath,
	}

	for _, item := range outputArr {
		result.Author = result.Author + " " + "@" + strings.TrimSpace(item)
	}

	return result
}

func generateGitOwnerFile(owners []LineResult, outputFileName string) string {
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	var message string
	for _, owner := range owners {
		message = message + fmt.Sprintf("%s %s", owner.FilePath, owner.Author) + "\n"
	}

	err = ioutil.WriteFile(file.Name(), []byte(message), 0644)
	if err != nil {
		panic(err)
	}
	return file.Name()
}

//var descriptionRegex = "\\/\\*\\*"

func printGitResult(results []LineResult) {
	for _, result := range results {
		if result.Error != nil {
			continue
		}
		fmt.Printf("%s %s \n", result.FilePath, result.Author)
	}
}
