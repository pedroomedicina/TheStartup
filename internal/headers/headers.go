package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return 0, false, nil
	}
	if crlfIndex == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:crlfIndex], []byte(":"), 2)
	key := string(parts[0])
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	h.Set(key, string(value))
	return crlfIndex + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
