package headers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069                           \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 57, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n a bunch of extra data")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid header with capital letters in the key
	headers = NewHeaders()
	data = []byte("hOSt: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Capital letters in header key should be lowercased
	headers = NewHeaders()
	data = []byte("Content-Type: application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, 32, n)
	assert.False(t, done)

	// Test: Invalid character in header key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character in header key")
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeadersTokenCharacters(t *testing.T) {
	// Test: All valid token special characters
	validTokenHeaders := []string{
		"X-Custom!: value1",
		"Header#: value2",
		"Field$: value3",
		"Name%: value4",
		"Key&: value5",
		"Test': value6",
		"Data*: value7",
		"Info+: value8",
		"My-Header: value9",       // hyphen
		"File.Extension: value10", // dot
		"Power^: value11",
		"Under_score: value12", // underscore
		"Back`tick: value13",   // backtick
		"Pipe|: value14",       // pipe
		"Tilde~: value15",      // tilde
	}

	for i, headerLine := range validTokenHeaders {
		t.Run(fmt.Sprintf("ValidToken_%d", i), func(t *testing.T) {
			headers := NewHeaders()
			data := []byte(headerLine + "\r\n\r\n")
			n, done, err := headers.Parse(data)
			require.NoError(t, err, "Should accept valid token character in header: %s", headerLine)
			require.NotNil(t, headers)
			assert.False(t, done)
			assert.Greater(t, n, 0)
			assert.Len(t, headers, 1, "Should have exactly one header")
		})
	}

	// Test: Mixed alphanumeric and token characters
	headers := NewHeaders()
	data := []byte("X-Custom-Header_123.test: application/json\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["x-custom-header_123.test"])
	assert.False(t, done)
	assert.Greater(t, n, 0)
}

func TestHeadersInvalidCharacters(t *testing.T) {
	// Test cases for invalid characters
	invalidHeaders := []struct {
		name   string
		header string
		char   string
	}{
		{"Space", "Invalid Header: value", " "},
		{"Tab", "Invalid\tHeader: value", "\t"},
		{"Newline", "Invalid\nHeader: value", "\n"},
		{"CarriageReturn", "Invalid\rHeader: value", "\r"},
		{"Unicode", "Invalid©Header: value", "©"},
		{"ControlChar", "Invalid\x00Header: value", "\x00"},
		{"HighBitASCII", "Invalid\x80Header: value", "\x80"},
		{"Parentheses", "Invalid(Header): value", "("},
		{"AngleBracket", "Invalid<Header: value", "<"},
		{"Comma", "Invalid,Header: value", ","},
		{"Semicolon", "Invalid;Header: value", ";"},
		{"Equals", "Invalid=Header: value", "="},
		{"Question", "Invalid?Header: value", "?"},
		{"AtSign", "Invalid@Header: value", "@"},
		{"SquareBracket", "Invalid[Header: value", "["},
		{"Backslash", "Invalid\\Header: value", "\\"},
		{"CurlyBrace", "Invalid{Header: value", "{"},
		{"DoubleQuote", "Invalid\"Header: value", "\""},
	}

	for _, tc := range invalidHeaders {
		t.Run("Invalid_"+tc.name, func(t *testing.T) {
			headers := NewHeaders()
			data := []byte(tc.header + "\r\n\r\n")
			n, done, err := headers.Parse(data)
			require.Error(t, err, "Should reject invalid character '%s' in header", tc.char)
			assert.Contains(t, err.Error(), "invalid character in header key")
			assert.Equal(t, 0, n)
			assert.False(t, done)
		})
	}
}

func TestHeadersEdgeCases(t *testing.T) {
	// Test: Empty header key (should fail)
	headers := NewHeaders()
	data := []byte(": value\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Single character header keys
	singleCharHeaders := []string{"A", "z", "0", "9", "!", "#", "$", "%", "&", "'", "*", "+", "-", ".", "^", "_", "`", "|", "~"}
	for _, char := range singleCharHeaders {
		t.Run("SingleChar_"+char, func(t *testing.T) {
			headers := NewHeaders()
			data := []byte(char + ": value\r\n\r\n")
			n, done, err := headers.Parse(data)
			require.NoError(t, err, "Single character '%s' should be valid", char)
			require.NotNil(t, headers)
			assert.Equal(t, "value", headers[strings.ToLower(char)])
			assert.False(t, done)
			assert.Greater(t, n, 0)
		})
	}

	// Test: Very long header key with all valid characters
	longKey := "Very-Long-Header-Name-With-Many-Valid-Characters_123.test!#$%&'*+-.^_`|~"
	headers = NewHeaders()
	data = []byte(longKey + ": value\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "value", headers[strings.ToLower(longKey)])
	assert.False(t, done)
	assert.Greater(t, n, 0)
}
