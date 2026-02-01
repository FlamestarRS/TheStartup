package request

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	reqLine, err := parseRequestLine(r)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(r []byte) (*RequestLine, error) {
	parts := strings.Split(string(r), "\r\n")
	subParts := strings.Split(parts[0], " ")
	if len(subParts) != 3 {
		return nil, fmt.Errorf("Malformed Request: %s", subParts)
	}
	method := subParts[0]
	requestTarget := subParts[1]
	httpVer := subParts[2]

	acceptableMethods := []string{"GET", "POST", "PUT", "DELETE"}

	if !slices.Contains(acceptableMethods, method) {
		return nil, fmt.Errorf("Malformed Method: %s", subParts)
	}

	if !strings.HasPrefix(requestTarget, "/") {
		return nil, fmt.Errorf("Malformed RequestTarget: %s", subParts)
	}

	if httpVer != "HTTP/1.1" {
		return nil, fmt.Errorf("Malformed HttpVersion: %s", subParts)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   strings.Trim(httpVer, "HTTP/"),
	}, nil
}
