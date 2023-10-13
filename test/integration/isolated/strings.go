package isolated

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// RandomName generates a random name of n length with the provided
// prefix. If prefix is omitted, the then entire name is random char.
func RandomName(prefix string, n int) string {
	if n == 0 {
		n = 32
	}
	if len(prefix) >= n {
		return prefix
	}
	p := make([]byte, n)
	_, _ = rand.Read(p)
	if len(prefix) > 0 {
		return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(p))[:n]
	}
	return hex.EncodeToString(p)[:n]
}
