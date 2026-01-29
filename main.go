package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	linesChan := getLinesChannel(f)

	for i := range linesChan {
		fmt.Printf("read: %s\n", i)
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
