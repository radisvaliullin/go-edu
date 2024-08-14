package httpserverv1

import (
	"bytes"
	"errors"
	"strconv"
)

const (
	HeaderFieldHost        = "Host"
	HeaderFieldContentType = "Content-Type"
	HeaderFieldContentLen  = "Content-Length"
)

func httpMessageHeaderDecoder(headerBytes []byte) (HttpHeader, error) {

	header := HttpHeader{}

	lines := bytes.Split(headerBytes, []byte{'\n'})
	if len(lines) < 1 {
		return header, errors.New("header decoder: header should has at least one line")
	}

	requestLine := lines[0]
	requestLineSegments := bytes.Split(requestLine, []byte{' '})
	if len(requestLineSegments) < 3 {
		return header, errors.New("header decoder: wrong request line segments")
	}
	header.Method = string(requestLineSegments[0])
	header.Path = string(requestLineSegments[1])

	// handle headers fields
	for _, line := range lines[1:] {
		sepIdx := bytes.Index(line, []byte{':'})
		if sepIdx == -1 {
			return header, errors.New("header decoder: wrong header fields separator")
		}
		fieldName := string(line[:sepIdx])
		fieldValue := string(bytes.TrimSpace(line[sepIdx+1:]))
		switch fieldName {
		case HeaderFieldHost:
			header.Host = fieldValue
		case HeaderFieldContentType:
			header.ContentType = fieldValue
		case HeaderFieldContentLen:
			l, err := strconv.Atoi(fieldValue)
			if err != nil {
				return header, errors.New("header decoder: wrong content length")
			}
			header.ContentLen = l
		}
	}

	return header, nil
}
