package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PallaviKapila/E-Commerce-Project/database"
	helper "github.com/PallaviKapila/E-Commerce-Project/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	//to hash our password

	"github.com/PallaviKapila/E-Commerce-Project/models"
	"github.com/gin-gonic/gin"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// in a database we can't store the password of the user so we have to hash it bcoz if somebody has
// access to the database then can't just take away all those passwords and start using them so we have to
// hash the password before storing it onto the database
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// boolean is going to store in passwordIsValid and string in msg
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	//we'll have one defition and we'll use bcrypt library to comapre hash and password
	//this function takes the provided password and also takes the other password, and compare both of them
	err := bcrypt.CompareHashAndPaasword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Password is incorrect")
		check = false
	}
	//check is boolean that we're returning
	return check, msg
}

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

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while cjecking for phone number"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exist"})
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
		//we worked on user id
		user.ID = primitive.NewObjectID()
		user.User_id = user.Id.Hex()
		//we want token, we are sending it to GenerateAllTokens function, this user_id is same that we created here
		token, refereshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		//here we have to set user.Token with token that i receive here
		user.Token = &token
		user.Refresh_token = &refereshToken
		//user object is created properly

		//now we need to insert this in database
		//we'll return as result as an insertion number
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			//s.printf to format the string
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		//if everything is ok
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		//the person who's trying to login that person needs to exist in our database
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//then we'll try and find that user, we'll try to find user from user collenction and we will use users' email
		//the found user we'll store in found Users after decoding it
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is incorrect"})
			return
		}
		//the user that is trying to login has sent his email and password to us, using that email we find that user and we store that in foundUser
		//then we use the user's password and found users's password to send it to a function verify password to check
		//if they're both matching
		passwordIsValid, msg := VerifyPassword(*user.Password, *&foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//one more thing we can add to check, one more validation
		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User Not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		//we'll use decode function and then we'll use foundUser to structure it
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//otherwise we'll say everrythig is alright
		c.JSON(http.StatusOK, foundUser)
	}
}

// gin is building a layer on top of http
// this fuunction can only be acess by the admin
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		//in authHelper we already have this function
		//we're just checking if user is ADMIN
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background())

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		//pipeline stage function in mongodb
		//or aggregration
		//for charts and stats
		matchStage := bson.D{{"$match", bson.D{{}}}}
		//all group functions start with something called id, id basically tells on which id u want data to be grouped on
		//group all ddata based on database and also help us create total count
		//$sum is not a pipeline function, it gets us calculate the sum of all the records basically
		//$push if we don't push everything to the root then we're not able to see the data we'll only see the count
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		//it helps us to define what all data points do we want to go to the user and which all shouldn't
		//it makes the data bit more readabale
		//
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		//return slice
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}
		//everything went well , we'll send all users to postman
		c.JSON(http.StatusOK, allusers[0])
	}
}

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
