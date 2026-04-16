package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	rr("Error setting up TCP listener: ", err)
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		rr("Error accepting connection: ", err)
		fmt.Println("Connection accepted!")
		ch := getLinesChannel(connection)
		for line := range ch {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Println("Connection closed!")

	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	slice := make([]byte, 8)
	var str string
	var sections []string

	go func() {
		for {
			n, err := f.Read(slice)
			if n == 0 {
				ch <- str
				close(ch)
				f.Close()
				return
			}
			rr("Error reading from reader: ", err)

			str += string(slice[:n])
			sections = strings.Split(str, "\n")
			if len(sections) == 2 {
				ch <- sections[0]

				str = sections[1]
			}

		}
	}()
	return ch
}

func rr(message string, err error) {
	if err != nil {
		fmt.Print("Error: ", err, "\n")
	}
}
