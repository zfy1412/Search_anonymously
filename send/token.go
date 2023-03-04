package send

import (
	"errors"
	"github.com/dgrijalva/jwt-go"

	"time"
)

// TokenExpireDuration is token's valid duration.
const TokenExpireDuration = time.Hour * 2

// EncryptionKey used for encryption.
var EncryptionKey = []byte("20180212")

// Claims is my custom claims.
type Claims struct {
	User
	jwt.StandardClaims
}

// GenerateToken is a function which generates token.
func GenerateToken(user User) (string, error) {
	claims := Claims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "yst",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(EncryptionKey)
}

// ParseToken is a function which parse token.
func ParseToken(key string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(key, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return EncryptionKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token not found or invalid")
}
