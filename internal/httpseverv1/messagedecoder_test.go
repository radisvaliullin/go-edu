package httpserverv1

import "testing"

func TestMessageDecoder(t *testing.T) {

	rawHttpRequestHeader := `POST /ping HTTP/1.1
Host: localhost:7373
User-Agent: curl/8.7.1
Accept: */*
Content-Type: application/json
Content-Length: 12`

	header, err := httpMessageHeaderDecoder([]byte(rawHttpRequestHeader))
	if err != nil {
		t.Fatalf("message header decode errro: %v", err)
	}
	t.Logf("message header decoded: %+v", header)
}
