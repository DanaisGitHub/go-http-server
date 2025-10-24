package request

import (
	"fmt"
	"io"
	"strings"
)

const (
	initialized int = iota
	done
)

type Request struct {
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string // 1.1 vs 2 vs ..
	RequestTarget string //?
	Method        string
}

const SEPARATOR string = "\r\n"

// GET /coffee HTTP/1.1
// /r/n

func parseRequestLine(s string) (*RequestLine, string, error) {
	reqlI := strings.Index(s, SEPARATOR)
	if reqlI == -1 {
		return nil, "", fmt.Errorf("no SEPARATOR = %v found", SEPARATOR)
	}
	requestLine := s[:reqlI]
	restOfReq := s[reqlI+len(SEPARATOR)-1:]

	// Seperate the Request line into 3 parts
	splitReqLine := strings.Split(requestLine, " ")
	if len(splitReqLine) != 3 {
		return nil, restOfReq, fmt.Errorf("malformed request line spaces")
	}

	httpVersion := strings.Split(splitReqLine[2], "/")
	if len(httpVersion) != 2 {
		return nil, restOfReq, fmt.Errorf("malformed http version")
	}
	if httpVersion[1] != "1.1" {
		return nil, restOfReq, fmt.Errorf("wrong http version")

	}

	rl := &RequestLine{
		HttpVersion:   httpVersion[1],
		Method:        splitReqLine[0],
		RequestTarget: splitReqLine[1],
	}

	return rl, restOfReq, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{}
	chunk := make([]byte, 8)
	streamBytes := make([]byte, 0)
	bytesRead := 0
	pos := 0
	r.state = initialized
	for {
		reqI, err := reader.Read(chunk)
		bytesRead += reqI
		streamBytes = append(streamBytes, chunk[:reqI]...)
		if err != nil && err == io.EOF {
			if err != nil {
				return nil, fmt.Errorf("EOF: %w", err)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk:=> %w", err)
		}
		num, err := r.parse(streamBytes)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse chunk: %w", err)
		}
		if r.state == done {
			pos = num
			break
		}

	}

	firstLine := streamBytes[:pos]

	rL, _, err := parseRequestLine(string(firstLine))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	r.RequestLine = *rL

	return r, nil

}

// finds the position of sepator and return in i index including length of SEPARATOR
func (r *Request) parse(data []byte) (int, error) {

	reqlI := strings.Index(string(data), SEPARATOR)
	if reqlI == -1 {
		return 0, nil
	}
	r.state = done
	return reqlI + len(SEPARATOR), nil

}
