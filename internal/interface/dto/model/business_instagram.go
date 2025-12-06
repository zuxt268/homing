package model

import (
	"time"
)

type BusinessInstagram struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string    `gorm:"column:name"`
	Memo          string    `gorm:"column:memo"`
	InstagramID   string    `gorm:"column:instagram_id"`
	InstagramName string    `gorm:"column:instagram_name"`
	BusinessName  string    `gorm:"column:business_name"`
	BusinessTitle string    `gorm:"column:business_title"`
	StartDate     time.Time `gorm:"column:start_date"`
	Status        int       `gorm:"column:status"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (*BusinessInstagram) TableName() string {
	return "business_instagrams"
}