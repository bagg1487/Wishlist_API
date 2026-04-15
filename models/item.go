package models

import "time"

type Item struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	WishlistID  uint   `gorm:"not null;index" json:"wishlist_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	URL         string `json:"url,omitempty"`
	Priority    int    `gorm:"default:1" json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsBooked bool `json:"is_booked" gorm:"-"`
}

type Booking struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	ItemID   uint      `gorm:"unique;not null;index" json:"item_id"`
	BookedBy string    `gorm:"not null" json:"booked_by"`
	BookedAt time.Time `json:"booked_at"`
}


