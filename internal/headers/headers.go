package headers

import (
	"bytes"
	"fmt"
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
	header := bytes.SplitN(data[:index], []byte(":"), 2)
	key := string(header[0])
	if strings.HasSuffix(key, " ") {
		return 0, false, fmt.Errorf("Error: field-name has trailing whitespace")
	}
	value := string(header[1])
	numBytes := index + 2
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(strings.Trim(value, "\r\n"))

	h[key] = value
	return numBytes, false, nil
}
