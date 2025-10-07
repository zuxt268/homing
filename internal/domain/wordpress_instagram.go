package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type WordpressInstagram struct {
	ID                 int
	Name               string
	WordpressDomain    string
	WordpressSiteTitle string
	InstagramID        string
	InstagramName      string
	Memo               string
	StartDate          time.Time
	Status             Status
	DeleteHash         bool
	CustomerType       CustomerType
	UpdatedAt          time.Time
	CreatedAt          time.Time
}

func (c *WordpressInstagram) GenerateAPIKey(secretPhrase string) string {
	data := []byte(secretPhrase + c.WordpressDomain)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

type Status int

type CustomerType int
