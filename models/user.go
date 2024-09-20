package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

var JwtKey = []byte("your_secret_key") // Replace with your actual secret key

// Claims defines the structure of the JWT payload
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Use HS256 instead of ES256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	signedToken, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
