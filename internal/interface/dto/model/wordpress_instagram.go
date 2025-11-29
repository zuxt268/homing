package model

import (
	"time"
)

type WordpressInstagram struct {
	ID                 int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name               string    `gorm:"column:name"`
	WordpressDomain    string    `gorm:"column:wordpress_domain"`
	WordpressSiteTitle string    `gorm:"column:wordpress_site_title"`
	InstagramID        string    `gorm:"column:instagram_id"`
	InstagramName      string    `gorm:"column:instagram_name"`
	Memo               string    `gorm:"column:memo"`
	StartDate          time.Time `gorm:"column:start_date"`
	Status             int       `gorm:"column:status"`
	DeleteHash         bool      `gorm:"column:delete_hash"`
	Categories         string    `gorm:"column:categories"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (*WordpressInstagram) TableName() string {
	return "wordpress_instagrams"
}
