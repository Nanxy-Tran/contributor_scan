package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

var ignoreFiles = []string{".git", ".idea"}

func main() {
	scanFolder("./")
}

func scanFolder(root string) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		var filePath string
		if root == "./" {
			filePath = root + file.Name()
		} else {
			filePath = root + "/" + file.Name()
		}

		if sliceContain(file.Name(), ignoreFiles) {
			continue
		}

		if file.IsDir() {
			scanFolder(filePath)
			continue
		}

		checkAuthor(filePath)
	}
}

func checkAuthor(filePath string) {
	cmd := exec.Command("bash", "-c", "git blame "+filePath+" --porcelain | grep '^author '")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output), filePath)

}

func sliceContain(searchValue string, slice []string) bool {
	for _, item := range slice {
		if item == searchValue {
			return true
		}
	}
	return false
}
