package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Customer struct {
	ID                           int
	Name                         string
	WordpressUrl                 string
	FacebookToken                string
	Email                        string
	Password                     string
	StartDate                    *time.Time
	InstagramBusinessAccountID   []string
	InstagramBusinessAccountName *string
	InstagramTokenStatus         int
	DeleteHash                   bool
	PaymentType                  string
	Type                         int
}

func (c *Customer) GenerateAPIKey(secretPhrase string) string {
	data := []byte(secretPhrase + c.WordpressUrl)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
