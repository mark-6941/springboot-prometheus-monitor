package protocol

func IsTLSPort(port uint16) bool {
	return port == 443 || port == 8443
}
