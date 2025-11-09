package headers

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
)

var SEPARATOR = []byte("\r\n")

type Headers map[string]string

// input I think will look like
// header:
	// field:value
	// field:value
	// ...
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// need to get a full line of the data

	n = len(data)
	done = false
	err = nil

	idx := bytes.Index(data, SEPARATOR)

	if idx == 0 { // no header
		return n, done, nil
	}
	if idx == len(data)-len(SEPARATOR)-1 {
		done = true
	}

	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("malformed header")
	}
	field := parts[0]
	value:=bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(field, []byte(" ")) {
		return 0, false, fmt.Errorf("field of header malformed, space(s) between field & :")

	}

	header:=[]string{
		string(field),string(value),
	}

	maps.Insert(h,slices.All(header))

	return n, done, err

}

func NewHeaders() Headers {
	return Headers{}
}
