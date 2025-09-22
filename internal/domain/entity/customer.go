package entity

import (
	"crypto/sha256"
	"encoding/hex"
)

type Customer struct {
	ID                 int
	WordpressUrl       string
	AccessToken        string
	InstagramAccountID string
}

func (c *Customer) GenerateAPIKey(secretPhrase string) string {
	data := []byte(secretPhrase + c.WordpressUrl)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
