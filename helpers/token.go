package helper

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	db "github.com/golang-jwtAuth-project/databases"
)

type SignedDetail struct {
	Id         string
	First_name string
	Last_name  string
	Email      string
	Phone      string
	User_type  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = db.Collection(db.Client, "user")

var secretKey string = "122398173981237198"

func GenerateToken(id, first, last, email, phone, userType string) (string, string, error) {

	claims := SignedDetail{
		Id:         id,
		First_name: first,
		Last_name:  last,
		Email:      email,
		Phone:      phone,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	freshClaims := SignedDetail{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(7*24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	freshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, freshClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return token, freshToken, nil
}

func UpdateToken(freshToken, token, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", token})
	updateObj = append(updateObj, bson.E{"fresh_token", freshToken})

	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": ObjId},
		bson.D{{"$set", updateObj}},
		&opts,
	)
	if err != nil {
		return err
	}
	return nil
}

func ValidateToken(singedToken string) (*SignedDetail, error) {

	token, err := jwt.ParseWithClaims(
		singedToken,
		&SignedDetail{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetail)
	if !ok {
		return nil, errors.New("Token is invalid!")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("Token is expired")
	}

	return claims, nil

}
