package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func checkAuthor(filePath string) <-chan LineResult {
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

//func printContributors(contributions []sortStruct) {
//	for _, contributor := range contributions {
//		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
//	}
//}
