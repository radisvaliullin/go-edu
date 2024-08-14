package httpserverv2

// wrap message with package prefix
func LogMsg(msg string) string {
	return "httpserver.v2: " + msg
}
