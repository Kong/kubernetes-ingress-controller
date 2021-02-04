package deckgen

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/kong/deck/file"
)

func GenerateSHA(targetContent *file.Content,
	customEntities []byte) ([]byte, error) {

	var buffer bytes.Buffer

	jsonConfig, err := json.Marshal(targetContent)
	if err != nil {
		return nil, fmt.Errorf("marshaling Kong declarative configuration to JSON: %w", err)
	}
	buffer.Write(jsonConfig)

	if customEntities != nil {
		buffer.Write(customEntities)
	}

	shaSum := sha256.Sum256(buffer.Bytes())
	return shaSum[:], nil
}
