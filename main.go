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
	defer f.Close()

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
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
		}
		currentLine += split[len(split)-1]
	}
	if currentLine != "" {
		fmt.Printf("read: %s\n", currentLine)
	}
}
