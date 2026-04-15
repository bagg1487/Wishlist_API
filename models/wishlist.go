package models

import "time"

type Wishlist struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	EventDate   *time.Time `json:"event_date,omitempty"`
	PublicToken string     `gorm:"unique;not null" json:"public_token"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}