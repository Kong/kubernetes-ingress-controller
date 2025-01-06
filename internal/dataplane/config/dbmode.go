package config

import "fmt"

type DBMode string

const (
	DBModeOff      DBMode = "off"
	DBModePostgres DBMode = "postgres"
)

func NewDBMode(mode string) (DBMode, error) {
	switch mode {
	case "", string(DBModeOff):
		return DBModeOff, nil
	case string(DBModePostgres):
		return DBModePostgres, nil
	}
	return "", fmt.Errorf("unsupported db mode: %q", mode)
}

// IsDBLessMode can be used to detect the proxy mode (db or dbless).
func (m DBMode) IsDBLessMode() bool {
	return m == "" || m == DBModeOff
}

// IsDBBacked returns true if the gateway is DB backed.
// reverse of IsDBLessMode for readability.
func (m DBMode) IsDBBacked() bool {
	return !m.IsDBLessMode()
}
