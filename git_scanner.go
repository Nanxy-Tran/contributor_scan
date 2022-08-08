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
	Author    string
	RawString string
	Error     error
}

type GitParser struct {
	method string
}

func (parser *GitParser) execute(fileChanel <-chan string) {
	if parser.method == "owner" {
		lineChannel := make(chan (<-chan LineResult), 10)
		for path := range fileChanel {
			fmt.Println(path)
			lineChannel <- checkAuthor(path)
		}

		close(lineChannel)

		var result []LineResult

		for lineChan := range lineChannel {
			select {
			case line := <-lineChan:
				fmt.Println("Line: ", line)
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
		fmt.Println("get ger")
		printGitResult(result)
		return
	}

	if parser.method == "contribute" {
		//TODO:
	}
}

func checkAuthor(filePath string) <-chan LineResult {
	authorChan := make(chan LineResult, 50)

	go func() {
		defer close(authorChan)
		cmd := exec.Command("bash", "-c", "git blame "+filePath+" --porcelain | sed 's/author //p' | sort | uniq -c| sort -rn | head -n 1 | sed 's/[0-9]*//g'")
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("Git blaming.... ", filePath)

		authorChan <- parseAuthor(string(output), filePath)
	}()

	return authorChan
}

func parseAuthor(output string, filePath string) LineResult {
	var outputArr = strings.Split(output, " ")
	var result = LineResult{
		FilePath: filePath,
	}

	for _, item := range outputArr {
		if item != "" {
			result.Author = item
			return result
		}
	}

	result.Author = ""
	return result
}

func generateGitOwnerFile(owners []LineResult, outputFileName string) string {
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	var message string
	for _, owner := range owners {
		message = message + fmt.Sprintf("%s @%s", owner.FilePath, owner.Author) + "\n"
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
		fmt.Printf("%s @%s", result.FilePath, result.Author)
	}
}
