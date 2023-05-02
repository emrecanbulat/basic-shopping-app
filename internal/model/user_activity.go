package model

import (
	"shoppingApp/internal/client"
	"time"
)

type UserActivity struct {
	ID          int64     `json:"id" gorm:"primarykey"` // Unique integer ID for the product
	UserID      int64     `json:"-"`
	UserName    string    `json:"-"`
	UserEmail   string    `json:"-"`
	User        User      `json:"user" gorm:"references:ID"`
	Type        string    `json:"type" gorm:"type:varchar(32);"`
	Action      string    `json:"action" gorm:"type:varchar(100);"`
	Description string    `json:"description" gorm:"type:text;"`
	CreatedAt   time.Time `json:"-"`
}

func (activity UserActivity) Create() (UserActivity, error) {
	result := client.PostgreSqlClient.Create(&activity)
	return activity, result.Error
}
