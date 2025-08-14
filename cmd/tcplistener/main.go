package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("error listening: %v\n", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println("\nShutting down listener...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v", err)
			break
		}

		fmt.Printf("New connection from %v\n", conn.RemoteAddr())
		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("%s\n", line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)
	go func() {
		currentLine := ""
		for {
			bytes := make([]byte, 8)
			numBytes, err := f.Read(bytes)
			if err != nil && err != io.EOF {
				fmt.Errorf("unexpected error reading file: %v", err)
			}

			if numBytes == 0 && err == io.EOF {
				break
			}

			part := strings.Split(string(bytes), "\n")
			currentLine += part[0]
			if len(part) > 1 {
				channel <- currentLine
				currentLine = part[1]
			}
		}

		channel <- currentLine
		close(channel)
		fmt.Println("Closing channel")
	}()

	return channel
}
