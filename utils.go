package main

import (
	"strings"
)

type sortStruct struct {
	Key   string
	Value int
}

//func (array Contributions) sort() []sortStruct {
//
//	var collection []sortStruct
//	for key, value := range array {
//		collection = append(collection, sortStruct{key, value})
//	}
//
//	sort.Slice(collection, func(i, j int) bool {
//		return collection[i].Value > collection[j].Value
//	})
//
//	return collection
//}

func filter(arr []string, key string) []string {
	var result []string
	for _, item := range arr {
		if strings.Contains(item, key) {
			continue
		} else {
			result = append(result, item)
		}
	}
	return result
}

func remove_space_arr(arr []string) []string {
	var result []string
	for _, item := range arr {
		if item == " " || item == "" {
			continue
		} else {
			result = append(result, item)
		}
	}
	return result
}
