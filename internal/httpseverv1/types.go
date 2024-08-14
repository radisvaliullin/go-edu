package httpserverv1

type HttpMessage struct {
	HttpHeader
	Payload []byte
}

type HttpHeader struct {
	Method string
	Path   string

	Host        string
	ContentType string
	ContentLen  int
}
