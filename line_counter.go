package main

const LINE_COUNTER = "LINE_COUNTER"

func countLines(fileChannel <-chan string) int {
	count := 0
	for range fileChannel {
		count++
	}
	return count
}
