package httpserverv1

import (
	"bufio"
	"bytes"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/radisvaliullin/go-edu/internal/utils/uerr"
)

type Config struct {
	Addr string
}

// HttpServer implements concurrent simple HTTP 1.1 Server using stdlib/net package.
// Just basic functionality.
// Support only /ping request for other return BAD REQUEST.
type HttpServer struct {
	config Config

	logger *slog.Logger
}

// New constructs HttpServer object.
// New is not reserved word, but it is commont to use it as constructor function.
// If constructed object's name same as package name then use just New,
// otherwise use NewObject.
func New(config Config, logger *slog.Logger) *HttpServer {
	s := &HttpServer{
		config: config,
		logger: logger,
	}
	return s
}

func (s *HttpServer) Start() error {

	ln, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.logger.Error(LogMsg("get linstener"), uerr.Error(err))
		return err
	}

	// accept connect from clients and handle concurrently
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error(LogMsg("accept connect"), uerr.Error(err))
			return err
		}
		go handleConnection(s.logger, conn, 1)
	}
}

func handleConnection(logger *slog.Logger, conn net.Conn, mode int) {
	defer conn.Close()

	// set conn read timeout (otherwise you can locked forever on read operation,
	// tcp do not guaranty connection liveness)
	if err := conn.SetReadDeadline(time.Now().Add(time.Second * 30)); err != nil {
		logger.Error(LogMsg("set conn read deadline"), uerr.Error(err))
		return
	}

	// read http message
	var httpMessage HttpMessage
	var err error
	// implemented two different methods to read http message from tcp stream
	// one use low level operations with bytes
	// second one use buffer reader
	switch mode {
	case 0:
		httpMessage, err = readHttpMessageLowLevel(logger, conn)
	case 1:
		httpMessage, err = readHttpMessage(logger, conn)
	default:
		httpMessage, err = readHttpMessage(logger, conn)
	}
	if err != nil {
		logger.Error(LogMsg("read http message"), uerr.Error(err))
		return
	}

	// send proper response
	var writeErr error
	switch httpMessage.Path {
	case "/ping":
		writeErr = pingResponse(conn)
	default:
		writeErr = errorResponse(conn)
	}
	if writeErr != nil {
		logger.Error(LogMsg("write response"), uerr.Error(writeErr))
	}
}

// uses bufio for read bytes
func readHttpMessage(logger *slog.Logger, reader io.Reader) (HttpMessage, error) {

	httpMessage := HttpMessage{}

	// read http header

	// limit http message size to 2048 bytes
	bufSize := 2048
	// buffered reader
	bufReader := bufio.NewReaderSize(reader, bufSize)

	header, err := httpMessageHeaderDecoder(bufReader)
	if err != nil {
		logger.Error(LogMsg("decode http header"), uerr.Error(err))
		return httpMessage, err
	}
	httpMessage.HttpHeader = header

	// read payload
	if header.ContentLen > 0 {

		// allocate new slice for payload to give dealloate buffer
		httpMessage.Payload = make([]byte, header.ContentLen)
		_, err := io.ReadFull(bufReader, httpMessage.Payload)
		if err != nil {
			logger.Error(LogMsg("read paylaod"), uerr.Error(err))
			return httpMessage, err
		}
	}

	return httpMessage, nil
}

// uses low level byte read operations
func readHttpMessageLowLevel(logger *slog.Logger, reader io.Reader) (HttpMessage, error) {

	httpMessage := HttpMessage{}

	// read http header
	// http header end with "\r\n\r\n"
	headerEndSep := []byte("\r\n\r\n")

	// limit http message size to 2048 bytes
	bufSize := 2048
	buf := make([]byte, bufSize)

	// read header
	readEndIdx := 0
	headerEndIdx := 0
	for {
		n, err := reader.Read(buf[readEndIdx:])
		if err != nil {
			logger.Error(LogMsg("read http header"), uerr.Error(err))
			return httpMessage, err
		}
		readEndIdx = readEndIdx + n
		if readEndIdx == bufSize {
			logger.Error(LogMsg("read buffer overflow"), uerr.Error(err))
			return httpMessage, err
		}
		readBytes := buf[:readEndIdx]

		// check that header was read
		headerEndIdx = bytes.Index(readBytes, headerEndSep)
		if headerEndIdx == -1 {
			continue
		}
		break
	}
	beginPayloadIdx := headerEndIdx + len(headerEndSep)

	// header decoding
	headerBytes := buf[:headerEndIdx]
	header, err := httpMessageHeaderDecoderLowLevel(headerBytes)
	if err != nil {
		logger.Error(LogMsg("decode http header"), uerr.Error(err))
		return httpMessage, err
	}
	httpMessage.HttpHeader = header

	// read payload
	endPayloadIdx := beginPayloadIdx + httpMessage.ContentLen
	if header.ContentLen > 0 {
		// need read until end of payload
		for {
			if readEndIdx >= endPayloadIdx {
				break
			}

			n, err := reader.Read(buf[readEndIdx:])
			if err != nil {
				logger.Error(LogMsg("read http payload"), uerr.Error(err))
				return httpMessage, err
			}
			readEndIdx = readEndIdx + n
			if readEndIdx == bufSize {
				logger.Error(LogMsg("read payload, buffer overflow"), uerr.Error(err))
				return httpMessage, err
			}

		}

		// allocate new slice for payload to give dealloate buffer
		httpMessage.Payload = make([]byte, header.ContentLen)
		copy(httpMessage.Payload, buf[beginPayloadIdx:endPayloadIdx])
	}

	return httpMessage, nil
}

func pingResponse(conn net.Conn) error {
	pingHttpResp := `HTTP/1.1 200 OK
Content-Length: 4
Content-Type: text/plain; charset=utf-8

pong`

	_, err := conn.Write([]byte(pingHttpResp))
	if err != nil {
		return err
	}
	return nil
}

func errorResponse(conn net.Conn) error {
	pingHttpResp := `HTTP/1.1 400 Bad Request
Content-Length: 15
Content-Type: text/plain; charset=utf-8

400 Bad Request`

	_, err := conn.Write([]byte(pingHttpResp))
	if err != nil {
		return err
	}
	return nil
}
