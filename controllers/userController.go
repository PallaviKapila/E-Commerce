package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/PallaviKapila/E-Commerce-Project/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-playground/validator/v10"
	//to hash our password

	helper "github.com/PallaviKapil/E-Commerce-Project/helpers"
	"github.com/gin-gonic/gin"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		//because signup function creates user in the database
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//handle validation, whatever data we have in user we're going to compare that out here and we're going to validate it, that's why we're using validate function
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
			//in postman we're returning count of the user
		}
		//here we'll use count to help us validate
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while cjecking for phone number"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exist"})
		}
	}
}

func Login()

func GetUsers()

// only admin can have access to users data, this means data coming from this api if it's data of another user then regular user can't access it only admin can access it
// gin gives access to its ownhandler function
func GetUser() gin.HandlerFunc {
	//
	return func(c *gin.Context) {
		userId := c.Param("*user_id")

		//I'll call my helper file and have function MatchUserTypeToUid to check user is admin or not
		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			//we can send bad request like this
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//which is already defined in userModel, we're creating this variable which is representation of this struct and now we can easily find the user from database using user_id
		var user models.User
		//we need to decode this data
		//we already knows that mongodb saves data in json and golang doesn't understand json that's why we created struct in first place
		//now we can decode the json information into the info. that golang understands which is string and all of those that's we used the decode function
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//after decoding we'll get the data in user we'll pass the data of user
		c.JSON(http.StatusOK, user)
	}
}
