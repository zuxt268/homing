package model

import "time"

type Token struct {
	ID       int       `gorm:"column:id;primaryKey;autoIncrement"`
	Token    string    `gorm:"column:token"`
	UpdateAt time.Time `gorm:"column:update_at;autoUpdateTime"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime"`
}

func (t *Token) TableName() string {
	return "token"
}
