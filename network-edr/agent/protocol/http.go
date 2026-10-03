package protocol

func IsHTTPPort(port uint16) bool {
	return port == 80 || port == 8080 || port == 8000 || port == 8081
}
