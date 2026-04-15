package utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"wishlist-api/models"
)

func GenerateJWT(userID uint) string {
	tk := &models.Token{UserID: userID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenStr, _ := token.SignedString([]byte(os.Getenv("token_pass")))
	return tokenStr
}