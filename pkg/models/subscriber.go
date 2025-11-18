package models

import "gorm.io/gorm"

// Subscriber represents a user subscribed to notifications
type Subscriber struct {
	gorm.Model
	ChatID    int64  `gorm:"uniqueIndex;not null" json:"chat_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}

// TableName specifies the table name for Subscriber model
func (Subscriber) TableName() string {
	return "subscribers"
}
