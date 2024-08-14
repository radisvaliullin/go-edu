package httpserverv1

import (
	"bufio"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestHttpMessageHeaderDecoderLowLevel(t *testing.T) {

	rawHttpRequestHeader := `POST /ping HTTP/1.1
Host: localhost:7373
User-Agent: curl/8.7.1
Accept: */*
Content-Type: application/json
Content-Length: 12`

	header, err := httpMessageHeaderDecoderLowLevel([]byte(rawHttpRequestHeader))
	if err != nil {
		t.Fatalf("message header decode low level errro: %v", err)
	}
	t.Logf("message header decoded low level: %+v", header)
}

func TestHttpMessageHeaderDecoder(t *testing.T) {

	rawHttpRequestMessage := `POST /ping HTTP/1.1
Host: localhost:7374
User-Agent: curl/8.7.1
Accept: */*
Content-Type: application/json
Content-Length: 24

`

	reader := strings.NewReader(rawHttpRequestMessage)
	bufReader := bufio.NewReaderSize(reader, 2048)

	header, err := httpMessageHeaderDecoder(bufReader)
	if err != nil {
		t.Fatalf("message header decode errro: %v", err)
	}
	t.Logf("message header decoded: %+v", header)
}

func TestReadHttpMessageLowLevel(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	rawHttpRequest := `POST /ping HTTP/1.1
Host: localhost:7373
User-Agent: curl/8.7.1
Accept: */*
Content-Type: application/json
Content-Length: 25`
	rawHttpRequest += "\r\n\r\n" + `{"abc":1234, "zxcv":2345}` + "\r\n"

	message, err := readHttpMessageLowLevel(logger, strings.NewReader(rawHttpRequest))
	if err != nil {
		t.Fatalf("message decode low level errro: %v", err)
	}
	t.Logf("message decoded low level: %+v", message)
}

func TestReadHttpMessage(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	rawHttpRequest := `POST /ping HTTP/1.1
Host: localhost:7373
User-Agent: curl/8.7.1
Accept: */*
Content-Type: application/json
Content-Length: 12

{"abc":1234}`

	message, err := readHttpMessage(logger, strings.NewReader(rawHttpRequest))
	if err != nil {
		t.Fatalf("message decode errro: %v", err)
	}
	t.Logf("message decoded: %+v", message)
}
