package httpserverv1

// wrap message with package prefix
func LogMsg(msg string) string {
	return "httpserver.v1: " + msg
}
