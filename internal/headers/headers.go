package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	index := bytes.Index(data, []byte("\r\n"))
	if index == -1 {
		return 0, false, nil
	}
	if index == 0 {
		return 2, true, nil
	}
	numBytes := index + 2
	header := bytes.SplitN(data[:index], []byte(":"), 2)
	key := strings.TrimLeft(string(header[0]), " ")
	validChars := regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+-.^_`|~]+$")
	if !validChars.MatchString(key) {
		return 0, false, fmt.Errorf("Error: field-name contains invalid characters: %s", key)
	}
	key = strings.ToLower(key)
	value := strings.TrimSpace(string(header[1]))

	if exists, ok := h[key]; ok {
		h[key] = exists + ", " + value
	} else {
		h[key] = value
	}

	return numBytes, false, nil
}
