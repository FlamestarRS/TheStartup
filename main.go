package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err == io.EOF {
			break
		}
		fmt.Printf("read: %s\n", data[:n])
	}
}
