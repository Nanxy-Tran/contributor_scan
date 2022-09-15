package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type LineResult struct {
	FilePath string
	Authors  []string
}

type GitParser struct {
	ignoreFiles []string
}

type GitProcessor struct{}

func (parser *GitParser) construct() {
	parser.ignoreFiles = []string{".gitignore", ".git", ".idea", ".jest", ".codeclimate.yml", "node_modules", "android/", "ios/", "coverage/", "png", "gif", "svg"}
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
	result := make(map[string][]string)
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
			result[parseFolderPath(line.FilePath)] = append(result[parseFolderPath(line.FilePath)], line.Authors...)
		}
	}

	for lineChan := range lineChans {
		select {
		case line := <-lineChan:
			result[parseFolderPath(line.FilePath)] = append(result[parseFolderPath(line.FilePath)], line.Authors...)
		default:
		}
	}

	generateGitOwnerOutput(result)
}

func checkAuthor(filePath string) <-chan LineResult {
	authorChan := make(chan LineResult, 50)

	go func() {
		defer close(authorChan)
		cmd := exec.Command("bash", "-c", "git blame "+filePath+" --line-porcelain |"+
			" grep '^author-mail' | sort | uniq -c |"+
			" sort -rn | head -n 2 | sed 's/[0-9]*//g' |"+
			" sed 's/author-mail <*//g' |"+
			" sed 's/>//g'")

		output, err := cmd.CombinedOutput()

		fmt.Println("Git blaming.... ", filePath)
		if err != nil {
			log.Fatal(err.Error())
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
		Authors:  []string{},
	}

	if len(outputArr) == 0 {
		return result
	}

	for _, author := range outputArr {
		if author == "" {
			continue
		}
		result.Authors = append(result.Authors, strings.TrimSpace(author))
	}

	return result
}

func parseFolderPath(path string) string {
	pathElements := strings.Split(path, "/")
	//Check is root folder
	if len(pathElements) > 2 {
		return strings.Join(removeLastElement(pathElements), "/")
	}
	return path
}

func generateGitOwnerOutput(owners map[string][]string) {
	var message string
	for path, authors := range owners {
		//TODO: customize how many code owners to pick
		mostContributedAuthors := getFrequentAuthor(authors)
		pickableOwners := 0

		if totalAuthors := len(mostContributedAuthors); totalAuthors > 2 {
			pickableOwners = 2
		} else {
			pickableOwners = totalAuthors
		}

		var authorIdentifierText = ""
		for _, author := range mostContributedAuthors[:pickableOwners] {
			authorIdentifierText += " " + author
		}
		message = message + fmt.Sprintf("%s %s", path, authorIdentifierText) + "\n"
	}
	fmt.Println(message)
}

func getFrequentAuthor(authors []string) []string {
	var result []string
	var authorSpace = make(map[string]int)

	for _, author := range authors {
		if _, ok := authorSpace[author]; ok {
			authorSpace[author]++
		} else {
			authorSpace[author] = 1
		}
	}

	for author := range authorSpace {
		result = append(result, author)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return authorSpace[result[i]] < authorSpace[result[j]]
	})

	return result
}

//var descriptionRegex = "\\/\\*\\*"

//func printContributors(contributions []sortStruct) {
//	for _, contributor := range contributions {
//		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
//	}
//}
