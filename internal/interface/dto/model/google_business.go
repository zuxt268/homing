package model

import (
	"time"
)

type GoogleBusiness struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name"`
	Title     string    `gorm:"column:title"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (*GoogleBusiness) TableName() string {
	return "google_businesses"
}