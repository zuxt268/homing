package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type WordpressInstagram struct {
	ID           int
	Name         string
	Wordpress    string
	InstagramID  string
	Memo         string
	StartDate    time.Time
	Status       Status
	DeleteHash   bool
	CustomerType CustomerType
}

func (c *WordpressInstagram) GenerateAPIKey(secretPhrase string) string {
	data := []byte(secretPhrase + c.Wordpress)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

type Status int

type CustomerType int
