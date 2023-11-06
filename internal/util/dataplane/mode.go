package dataplane

type DBMode string

const (
	DBModeOff       = "off"
	DBModePostgres  = "postgres"
	DBModeCassandra = "cassandra"
)

// IsDBLessMode can be used to detect the proxy mode (db or dbless).
func IsDBLessMode(mode DBMode) bool {
	return mode == "" || mode == DBModeOff
}

// DBBacked returns true if the gateway is DB backed.
// reverse of IsDBLessMode for readability.
func DBBacked(mode DBMode) bool {
	return !IsDBLessMode(mode)
}
