package main

type LineParser interface {
	execute(fileChannel <-chan string)
}

//func (scanner *RawStringCrawler) scan(filePath string) <-chan LineResult {
//	lineChan := make(chan LineResult, 50)
//
//	fileLines, err := readLines(filePath)
//
//	if err != nil {
//		defer close(lineChan)
//		return lineChan
//	}
//
//	go func() {
//		defer close(lineChan)
//		regex := regexp.MustCompile(descriptionRegex)
//		for _, line := range fileLines {
//			descriptionLine := regex.FindString(line)
//			if descriptionLine != "" {
//				lineChan <- LineResult{FilePath: filePath}
//				return
//			}
//		}
//
//		lineChan <- LineResult{FilePath: filePath}
//	}()
//
//	return lineChan
//}

// LineCrawlerFactory TODO: Init the predicate
func LineParserFactory(crawler string) LineParser {
	return &GitParser{method: "owner"}
}

func main() {
	//start := time.Now()
	//filesChan := make(chan string, 10)
	//var lineChans chan (<-chan LineResult)
	//lineChans = make(chan (<-chan LineResult), 50)

	crawler := FileCrawler{ignoreFiles: ignoreFiles}
	crawler.scan("./")
	parser := LineParserFactory("counter")
	parser.execute(crawler.fileChannels)

	//go scanFolder("./", filesChan)
	//
	//totalLines := countLines(filesChan)
	////totalLines := countLines(filesChan)
	//fmt.Println("Number of lines: ", totalLines)

	//go func() {
	//	defer close(lineChans)
	//	for path := range filesChan {
	//		fmt.Println("Checking.... ", path)
	//		lineChans <- checkLines(path)
	//	}
	//}()
	//
	//for lineChan := range lineChans {
	//	select {
	//	case line := <-lineChan:
	//		result = append(result, line)
	//	}
	//}
	//
	//for lineChan := range lineChans {
	//	select {
	//	case line := <-lineChan:
	//		result = append(result, line)
	//	default:
	//	}
	//}
	//
	//fmt.Println("Total typescript files: ", len(result))
	//var documentedFiles = filter(result, "NO DESCRIPTION FOUND")
	//fmt.Println("Number of documentation files: ", len(documentedFiles))
	//fmt.Println("Documentation coverage: ", math.Round(float64(len(documentedFiles))/float64(len(result))*100), " %")
	//fmt.Printf("Executed in %s", time.Since(start))
}

func appendValidLine(path string) {

}

var descriptionRegex = "\\/\\*\\*"

//func printContributors(contributions []sortStruct) {
//	for _, contributor := range contributions {
//		fmt.Println("Contributor: " + contributor.Key + " has contributed " + strconv.Itoa(contributor.Value) + " files")
//	}
//}
