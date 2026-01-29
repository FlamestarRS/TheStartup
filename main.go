package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		fmt.Println("Waiting for Conenction...")
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection Accepted...")
		linesChan := getLinesChannel(conn)

		for i := range linesChan {
			fmt.Println(i)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)
	go func() {
		defer f.Close()
		defer close(channel)
		currentLine := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err == io.EOF {
				break
			}

			split := strings.Split(string(data[:n]), "\n")

			for i := 0; i < len(split)-1; i++ {
				currentLine += split[i]
				channel <- currentLine
				currentLine = ""
			}
			currentLine += split[len(split)-1]
		}
		if currentLine != "" {
			channel <- currentLine
		}
	}()
	return channel
}
