package res

import "time"

type Token struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expired_at"`
}
