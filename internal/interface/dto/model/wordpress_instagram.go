package model

import (
	"time"
)

type WordpressInstagram struct {
	ID           int       `gorm:"column:id;primaryKey"`
	Name         string    `gorm:"column:name"`
	Wordpress    string    `gorm:"column:wordpress"`
	InstagramID  string    `gorm:"column:instagram_id"`
	Memo         string    `gorm:"column:memo"`
	StartDate    time.Time `gorm:"column:start_date"`
	Status       int       `gorm:"column:status"`
	DeleteHash   bool      `gorm:"column:delete_hash"`
	CustomerType int       `gorm:"column:customer_type"`
	UpdateAt     time.Time `gorm:"column:update_at;autoUpdateTime"`
	CreateAt     time.Time `gorm:"column:create_at;autoCreateTime"`
}

func (*WordpressInstagram) TableName() string {
	return "wordpress_instagrams"
}
