package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	files, err := ioutil.ReadDir("logs")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}
