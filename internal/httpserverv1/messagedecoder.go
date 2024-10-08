package httpserverv1

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
)

const (
	HeaderFieldHost        = "Host"
	HeaderFieldContentType = "Content-Type"
	HeaderFieldContentLen  = "Content-Length"
)

func httpMessageHeaderDecoder(reader *bufio.Reader) (HttpHeader, error) {

	header := HttpHeader{}

	// handle header first line
	requestLine, _, err := reader.ReadLine()
	if err != nil {
		return header, err
	}
	requestLineSegments := bytes.Split(requestLine, []byte{' '})
	if len(requestLineSegments) < 3 {
		return header, errors.New("header decoder: wrong request line segments")
	}
	header.Method = string(requestLineSegments[0])
	header.Path = string(requestLineSegments[1])

	// handle headers fields
	cnt := 0
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return header, err
		}
		// if empty line stop iterate
		if len(line) == 0 {
			cnt++
			break
		}

		sepIdx := bytes.Index(line, []byte{':'})
		if sepIdx == -1 {
			return header, errors.New("header decoder: header fields separator not found")
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

func httpMessageHeaderDecoderLowLevel(headerBytes []byte) (HttpHeader, error) {

	header := HttpHeader{}

	// split returns slices pointing to same underlying array
	lines := bytes.Split(headerBytes, []byte{'\n'})
	if len(lines) < 1 {
		return header, errors.New("header decoder low level: header should has at least one line")
	}

	// decode first line of message
	requestLine := lines[0]
	requestLineSegments := bytes.Split(requestLine, []byte{' '})
	if len(requestLineSegments) < 3 {
		return header, errors.New("header decoder low level: wrong request line segments")
	}
	header.Method = string(requestLineSegments[0])
	header.Path = string(requestLineSegments[1])

	// handle header fields
	for _, line := range lines[1:] {
		sepIdx := bytes.Index(line, []byte{':'})
		if sepIdx == -1 {
			return header, errors.New("header decoder low level: header fields separator not found")
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
				return header, errors.New("header decoder low level: wrong content length")
			}
			header.ContentLen = l
		}
	}

	return header, nil
}
