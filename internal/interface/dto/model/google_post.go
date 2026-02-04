package model

import (
	"time"
)

type GooglePost struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	InstagramURL string    `gorm:"column:instagram_url"`
	MediaID      string    `gorm:"column:media_id"`
	CustomerID   int       `gorm:"column:customer_id"`
	Name         string    `gorm:"column:name"`
	GoogleURL    string    `gorm:"column:google_url"`
	CreateTime   string    `gorm:"column:create_time"`
	PostType     string    `gorm:"column:post_type"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (*GooglePost) TableName() string {
	return "google_posts"
}