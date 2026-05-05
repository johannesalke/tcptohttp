package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:55555")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

	}

}
