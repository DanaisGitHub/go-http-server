package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var NEWLINE = []byte("\r\n")

type Headers map[string]string

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) Set(field, value string) {
	field = strings.ToLower(field)
	h[field] = value
}

func isToken(fieldName string) bool {
	return true
}

// this needs to be able to parse multiple lines
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// need to get a full line of the data

	n = 0
	done = false
	err = nil
	for {
		idx := bytes.Index(data[n:], NEWLINE)

		if idx == -1 {
			break
		}
		if idx == 0 { // no header
			done = true
			break
		}

		currData := data[n:idx]
		n += idx + len(NEWLINE)

		parts := bytes.SplitN(currData, []byte(":"), 2)
		if len(parts) != 2 {
			return 0, false, fmt.Errorf("malformed header")
		}
		field := parts[0]
		value := bytes.TrimSpace(parts[1])
		if bytes.HasSuffix(field, []byte(" ")) {
			return 0, false, fmt.Errorf("field of header malformed, space(s) between field & :")
		}

		if !isToken(string(field)) {
			return 0, false, fmt.Errorf("field name %s is not valid token", field)
		}

		h.Set(string(field), string(value))
	}
	return n, done, err
}

func NewHeaders() Headers {
	return Headers{}
}
