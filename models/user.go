package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username   string             `json:"username"`
	Password   string             `json:"password"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Role       string             `json:"role"`
	Favourites []Favourite        `json:"favourites"`
}

type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}
