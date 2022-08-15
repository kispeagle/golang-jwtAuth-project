package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id          primitive.ObjectID `bson:"_id"`
	First_name  string             `json:"first_name" validate:"required|min=2|max=100"`
	Last_name   string             `json:"last_name" validate:"required|min=2|max=100"`
	Account     string             `json:"account" validate:"required|min=4|max=100"`
	Email       string             `json:"email" validate:"required|email"`
	Password    string             `json:"password" validate:"required|min=9|max=200"`
	Phone       string             `json:"phone"`
	User_type   string             `json:"user_type" validate:"required|eq=ADMIN|eq=USER"`
	Token       string             `json:"token"`
	Fresh_token string             `json:"fresh_token"`
	Create_at   time.Time          `json:"create_at"`
	Modified_at time.Time          `json:"modified_at"`
}
