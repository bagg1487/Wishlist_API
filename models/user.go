package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"unique;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
}

type Token struct {
	UserID uint
	jwt.StandardClaims
}