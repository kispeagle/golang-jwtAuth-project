package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	db "github.com/golang-jwtAuth-project/databases"
	helper "github.com/golang-jwtAuth-project/helpers"
	model "github.com/golang-jwtAuth-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = db.Collection(db.Client, "user")
var validator_ = validator.New()

func hashPassword(password string) (string, error) {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(HashedPassword), err
}

func VerifyPassword(InputPassword, ProvidedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(InputPassword), []byte(ProvidedPassword))
	if err != nil {
		return err
	}
	return nil
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user model.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		err = validator_.Struct(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"account": user.Account})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"Error": "Account already existed"})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"Error": "Email already existed"})
			return
		}
		fmt.Println(count)

		count, err = userCollection.CountDocuments(ctx, bson.M{"Phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"Error": "Phone number already existed"})
			return
		}

		user.Id = primitive.NewObjectID()
		user.Create_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		log.Println(err)
		user.Modified_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		log.Println(err)
		user.Password, err = hashPassword(user.Password)
		log.Println(err)

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user model.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		if user.Password == "" || user.Account == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Lack of information"})
			return
		}

		var providedUser model.User
		err = userCollection.FindOne(ctx, bson.M{"account": user.Account}).Decode(&providedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		if err := VerifyPassword(providedUser.Password, user.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Passord is wrong"})
			return
		}

		token, freshToken, err := helper.GenerateToken(
			providedUser.Id.Hex(),
			providedUser.First_name,
			providedUser.Last_name,
			providedUser.Email,
			providedUser.Phone,
			providedUser.User_type,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		err = helper.UpdateToken(freshToken, token, providedUser.Id.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})

	}
}

func Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := helper.CheckRole(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
			return
		}

		var users []bson.M
		result, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		err = result.All(ctx, &users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)

	}
}

func GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		idString := c.Param("id")
		err := helper.MatchId(c, idString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
			return
		}

		objId, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		var user model.User
		err = userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		idHex := c.Param("id")
		err := helper.MatchId(c, idHex)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
			return
		}

		var updateObj primitive.D
		var updateInfo model.User
		err = c.BindJSON(&updateInfo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		if updateInfo.Email != "" {
			updateObj = append(updateObj, bson.E{"email", updateInfo.Email})
		}

		if updateInfo.Password != "" {
			password, _ := hashPassword(updateInfo.Password)
			updateObj = append(updateObj, bson.E{"password", password})
		}

		if updateInfo.First_name != "" {
			updateObj = append(updateObj, bson.E{"first_name", updateInfo.First_name})
		}

		if updateInfo.Last_name != "" {
			updateObj = append(updateObj, bson.E{"last_name", updateInfo.Last_name})
		}

		if updateInfo.Phone != "" {
			updateObj = append(updateObj, bson.E{"phone", updateInfo.Phone})
		}

		now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"modified_at", now})

		upsert := true
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}

		objId, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.D{{"$set", updateObj}}, &opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

		idHex := c.Param("id")
		objId, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		if err := helper.MatchId(c, idHex); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
			return
		}

		result, err := userCollection.DeleteOne(c, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}
