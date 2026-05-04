package main

import (
	"fmt"
	"strings"
)

const example = "this:is:one:example"

func main() {
	fmt.Println(strings.Join(strings.SplitN(example, ":", 2), "\n"))
	mapp := make(map[string]string)

	testResult := mapp["test"]
	if testResult == "" {
		fmt.Println("Correct!")
	}
	fmt.Print(testResult == "")
}
