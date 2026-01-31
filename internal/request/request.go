package request

import (
	"fmt"
	"io"
	"log"
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
		return &Request{}, err
	}

	req := &Request{
		RequestLine: reqLine,
	}

	return req, nil
}

func parseRequestLine(r []byte) (RequestLine, error) {
	parts := strings.Split(string(r), "\r\n")
	subParts := strings.Split(parts[0], " ")
	if len(subParts) != 3 {
		return RequestLine{}, fmt.Errorf("Malformed Request: %s", r)
	}

	requestLine := RequestLine{
		Method:        subParts[0],
		RequestTarget: subParts[1],
		HttpVersion:   subParts[2][5:],
	}

	return requestLine, nil
}
