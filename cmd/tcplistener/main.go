package main

import (
	"TheStartup/internal/request"
	"fmt"
	"log"
	"net"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("error listening: %v\n", err)
		return
	}

	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v", err)
			break
		}

		fmt.Printf("New connection from %v\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %v", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Printf("Headers:\n")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}
}
