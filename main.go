package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()

	cliArgs := os.Args[1:]

	var parserType = ""
	var executeMethod = ""

	if len(cliArgs) > 0 {
		parserType = cliArgs[0]
		if cliArgs[0] == "line" {
			executeMethod = cliArgs[0]
		} else {
			executeMethod = "owner"
		}
	}

	filesChan := make(chan string)
	parser := parserFactory(parserType)
	parser.construct()

	processor := fileProcessorFactory(executeMethod)
	go parser.scanFolder("./", filesChan)

	processor.execute(filesChan)
	fmt.Printf("Executed in %s \n", time.Since(start))
}
