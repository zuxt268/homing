package model

import "time"

type Customer struct {
	ID                           int       `gorm:"column:id;primaryKey"`
	Name                         string    `gorm:"column:name"`
	Email                        string    `gorm:"column:email"`
	Password                     string    `gorm:"column:password"`
	WordpressURL                 string    `gorm:"column:wordpress_url"`
	FacebookToken                string    `gorm:"column:facebook_token"`
	StartDate                    time.Time `gorm:"column:start_date"`
	InstagramBusinessAccountID   string    `gorm:"column:instagram_business_account_id"`
	InstagramBusinessAccountName string    `gorm:"column:instagram_business_account_name"`
	InstagramTokenStatus         int       `gorm:"column:instagram_token_status"`
	DeleteHash                   bool      `gorm:"column:delete_hash"`
	PaymentType                  string    `gorm:"column:payment_type"`
	Type                         int       `gorm:"column:type"`
}
