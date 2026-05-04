package main

import (
	"fmt"
	"github.com/johannesalke/tcptohttp/internal/headers"
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
	fmt.Println(testResult == "")
	h := headers.NewHeaders()
	h["test"] = "This"
	fmt.Println(h.Get("tEsT"))
	if h.Get("oops") == "" {
		fmt.Println("Still correct!")
	}
}
