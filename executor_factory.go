package main

type FileProcessor interface {
	execute(filesChan <-chan string)
}

func fileProcessorFactory(method string) FileProcessor {
	if method == "owner" {
		return &GitProcessor{}
	}
	return &LineProcessor{method: method}
}
