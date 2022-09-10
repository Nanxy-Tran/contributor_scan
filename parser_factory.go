package main

type Parser interface {
	construct()
	scanFolder(root string, fileChans chan<- string)
}

func parserFactory(parserType string) Parser {
	if parserType == "git" {
		return &GitParser{}
	}
	return &LineParser{}
}
