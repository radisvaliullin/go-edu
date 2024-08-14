package httpserverv1

import (
	"bytes"
	"log/slog"
	"net"
	"time"

	"github.com/radisvaliullin/go-edu/internal/utils/uerr"
)

type Config struct {
	Addr string
}

// HttpServer implements concurrent simple Http1.1 Server using stdlib/net package.
// Just bacic functionality.
// Support only /ping request for other return bad request.
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
		go handleConnection(s.logger, conn)
	}
}

func handleConnection(logger *slog.Logger, conn net.Conn) {
	defer conn.Close()

	// set conn read timeout
	if err := conn.SetReadDeadline(time.Now().Add(time.Second * 30)); err != nil {
		logger.Error(LogMsg("set conn read deadline"), uerr.Error(err))
		return
	}
	httpMessage, err := readHttpMessageLowLevel(logger, conn)
	if err != nil {
		logger.Error(LogMsg("read http message"), uerr.Error(err))
		return
	}

	var writeErr error
	switch httpMessage.Path {
	case "/ping":
		writeErr = pingResponse(conn)
	default:
		writeErr = errorResponse(conn)
	}
	if writeErr != nil {
		logger.Error(LogMsg("response write"), uerr.Error(writeErr))
	}
}

func readHttpMessageLowLevel(logger *slog.Logger, conn net.Conn) (HttpMessage, error) {

	httpMessage := HttpMessage{}

	// read http header
	// http header end with "\r\n\r\n"
	// will use low lever read byte operatios
	// true way use [bufio] stdlib package
	headerEndSep := []byte("\r\n\r\n")

	// define message read buffer slice of bytes
	// consider that http message not bigger than 2048 bytes
	buf := make([]byte, 2048)

	// read header
	readEndIdx := 0
	headerEndIdx := 0
	for {
		n, err := conn.Read(buf[readEndIdx:])
		if err != nil {
			logger.Error(LogMsg("read http header"), uerr.Error(err))
			return httpMessage, err
		}
		if len(buf[readEndIdx:]) == n {
			logger.Error(LogMsg("read buffer overflow"), uerr.Error(err))
			return httpMessage, err
		}
		readEndIdx = readEndIdx + n
		readBytes := buf[:readEndIdx]

		// check that header was read
		headerEndIdx = bytes.Index(readBytes, headerEndSep)
		if headerEndIdx == -1 {
			continue
		}
		break
	}
	beginPayloadIdx := headerEndIdx + len(headerEndSep)

	// header
	headerBytes := buf[:headerEndIdx]
	header, err := httpMessageHeaderDecoder(headerBytes)
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

			n, err := conn.Read(buf[readEndIdx:])
			if err != nil {
				logger.Error(LogMsg("read http payload"), uerr.Error(err))
				return httpMessage, err
			}
			if len(buf[readEndIdx:]) == n {
				logger.Error(LogMsg("read payload, buffer overflow"), uerr.Error(err))
				return httpMessage, err
			}
			readEndIdx = readEndIdx + n
		}

		httpMessage.Payload = buf[beginPayloadIdx:endPayloadIdx]
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
Content-Length: 0

`

	_, err := conn.Write([]byte(pingHttpResp))
	if err != nil {
		return err
	}
	return nil
}
