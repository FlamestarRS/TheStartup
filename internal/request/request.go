package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0

	req := &Request{
		RequestLine: RequestLine{},
		State:       requestStateInitialized,
	}

	for req.State != requestStateDone {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err == io.EOF {
			req.State = requestStateDone
			break
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func parseRequestLine(r []byte) (*RequestLine, int, error) {
	data := len(r)
	parts := strings.Split(string(r), "\r\n")
	if len(parts) == 1 {
		return nil, 0, nil
	}
	subParts := strings.Split(parts[0], " ")
	if len(subParts) != 3 {
		return nil, data, fmt.Errorf("Malformed Request: %s", subParts)
	}
	method := subParts[0]
	requestTarget := subParts[1]
	httpVer := subParts[2]

	acceptableMethods := []string{"GET", "POST", "PUT", "DELETE"}

	if !slices.Contains(acceptableMethods, method) {
		return nil, data, fmt.Errorf("Malformed Method: %s", subParts)
	}

	if !strings.HasPrefix(requestTarget, "/") {
		return nil, data, fmt.Errorf("Malformed RequestTarget: %s", subParts)
	}

	if httpVer != "HTTP/1.1" {
		return nil, data, fmt.Errorf("Malformed HttpVersion: %s", subParts)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   strings.Trim(httpVer, "HTTP/"),
	}, data, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case requestStateInitialized:
		reqLine, numBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numBytes == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.State = 1
		return numBytes, nil
	case requestStateDone:
		return 0, fmt.Errorf("Error: Trying to read data in a done state")
	default:
		return 0, fmt.Errorf("Error: Unknown state")
	}
}
