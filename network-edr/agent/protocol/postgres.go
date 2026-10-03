package protocol

func IsPostgresPort(port uint16) bool {
	return port == 5432
}
