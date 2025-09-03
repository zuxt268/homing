package domain

import "time"

type Post struct {
	ID            int       `gorm:"column:id;primaryKey"`
	MediaID       string    `gorm:"column:media_id"`
	CustomerID    int       `gorm:"column:customer_id"`
	Timestamp     string    `gorm:"column:timestamp"`
	MediaURL      string    `gorm:"column:media_url"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	Permalink     string    `gorm:"column:permalink"`
	WordpressLink string    `gorm:"column:wordpress_link"`
}
