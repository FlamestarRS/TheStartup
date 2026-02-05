package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"TheStartup/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	State       requestState
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0

	req := &Request{
		RequestLine: RequestLine{},
		State:       requestStateInitialized,
		Headers:     headers.NewHeaders(),
	}

	for req.State != requestStateDone {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != requestStateDone {
					return nil, fmt.Errorf("Error: Missing end of headers")
				}
				break
			}
			return nil, err
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
	totalBytesParsed := 0
	for r.State != requestStateDone {
		numBytes, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += numBytes
		if numBytes == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
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
		r.State = requestStateParsingHeaders
		return numBytes, nil
	case requestStateParsingHeaders:
		numBytes, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing header: %s", err)
		}
		if done {
			r.State = requestStateDone
		}
		return numBytes, nil
	case requestStateDone:
		return 0, fmt.Errorf("Error: Trying to read data in a done state")
	default:
		return 0, fmt.Errorf("Error: Unknown state")
	}
}
