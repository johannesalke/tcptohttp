package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	fmt.Print("I hope I get the job!")
	f, err := os.Open("messages.txt")
	rr(err)
	slice := make([]byte, 8)
	var str string
	var sections []string

	for {
		n, err := f.Read(slice)
		if n == 0 {
			fmt.Printf("read: %s\n", str)
			os.Exit(0)
		}
		rr(err)

		str += string(slice[:n])
		sections = strings.Split(str, "\n")
		if len(sections) == 2 {
			fmt.Printf("read: %s\n", sections[0])
			str = sections[1]
		}

	}
}

func rr(err error) {
	if err != nil {
		fmt.Print("Error: ", err, "\n")
	}
}
