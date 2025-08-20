package headers

import (
	"bytes"
	"fmt"
	"slices"
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

	if !isValidHeaderKey([]byte(key)) {
		return 0, false, fmt.Errorf("invalid character in header key: %s", key)
	}

	h.Set(key, string(value))
	return crlfIndex + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func isValidHeaderKey(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}
