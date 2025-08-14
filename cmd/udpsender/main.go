package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	updaddr, err := net.ResolveUDPAddr("udp4", "localhost:42069")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
	}

	udpconn, err := net.DialUDP("udp4", nil, updaddr)
	if err != nil {
		fmt.Println("Error connecting to UDP:", err)
	}

	defer udpconn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
		}

		_, err = udpconn.Write([]byte(input))
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
		}
	}
}
