package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type LineResult struct {
	FilePath string
	Author   string
	Output   string
}

type GitParser struct {
	ignoreFiles []string
}

type GitProcessor struct{}

func (parser *GitParser) construct() {
	parser.ignoreFiles = []string{".gitignore", ".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/"}
}

func (parser *GitParser) scanFolder(root string, fileChan chan<- string) {
	files, err := os.ReadDir(root)

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

		for _, ignore := range parser.ignoreFiles {
			if strings.Contains(filePath, ignore) {
				continue FILES_LOOP
			}
		}

		if file.IsDir() {
			parser.scanFolder(filePath, fileChan)
		} else {
			fileChan <- filePath
			continue FILES_LOOP
		}
	}

	if root == "./" {
		close(fileChan)
	}
}

func (*GitProcessor) execute(filesChan <-chan string) {
	var lineChans chan (<-chan LineResult)
	var result []LineResult
	lineChans = make(chan (<-chan LineResult), 50)

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

	generateGitOwnerFile(result, "git_owner.txt")
}

func checkAuthor(filePath string) <-chan LineResult {
	authorChan := make(chan LineResult, 50)

	go func() {
		defer close(authorChan)
		cmd := exec.Command("bash", "-c", "git blame "+filePath+" --line-porcelain | grep '^author ' | sort | uniq -c| sort -rn | head -n 2 | sed 's/[0-9]*//g' | sed 's/author*//g'")
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

func parseAuthor(output string, filePath string) LineResult {
	var outputArr = strings.Split(strings.TrimSpace(output), "\n")

	var result = LineResult{
		FilePath: filePath,
		Author:   "",
	}

	if len(outputArr) == 0 {
		return result
	}

	for _, author := range outputArr {
		if result.Author != "" {
			result.Author = result.Author + " " + fmt.Sprintf("@%s", strings.TrimSpace(author))
		} else {
			result.Author = fmt.Sprintf("@%s", strings.TrimSpace(author))
		}
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

	err = os.WriteFile(file.Name(), []byte(message), 0644)
	if err != nil {
		panic(err)
	}
	return file.Name()
}

//var descriptionRegex = "\\/\\*\\*"

//func printContributors(contributions []sortStruct) {
//	for _, contributor := range contributions {
//		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
//	}
//}
