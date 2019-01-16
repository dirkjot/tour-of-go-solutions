package main

import (
	"fmt"
	"strings"
	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	result := map[string]int{}
	fields := strings.Fields(s)
	for _, word := range fields {
		if prevCount, ok := result[word]; !ok {
			result[word] = 1
		} else {
			result[word] = prevCount + 1
		}
	}
	return result
}

func main() {
	wc.Test(WordCount)
}
