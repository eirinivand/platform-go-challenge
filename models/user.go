package models

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User User model
//
// swagger:model User
type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username   string             `json:"username"`
	Password   string             `json:"password,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at" `
	ModifiedAt time.Time          `json:"modified_at" bson:"modified_at" `
	Role       string             `json:"role,omitempty"`
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}
