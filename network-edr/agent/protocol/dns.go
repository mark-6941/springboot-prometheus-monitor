package protocol

func IsDNSPort(port uint16) bool {
	return port == 53
}
