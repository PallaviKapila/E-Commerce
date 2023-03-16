// this file will have GenerateAllTokens
package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PallaviKapila/E-Commerce-Project/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// it's struct , jwt uses hashing mechanishm to basically take the details that we give to give us token
// jwt token also consists of secret key usually, secret key is something that is with us that's secret
type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

// sending the name of the collection that we want to access
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedTokens string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		//this small email is what we receive inside this function and Email is part of signed details which we have to complete
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			//with every token we need to have expiry
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	//it is used to get a new token if you're regular token has expired
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	//SigningMethodHS256 algorithmn to create encrypted token for you
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	//let's handle the token
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err

}

// everytime we login we'll get a new token and refresh_token
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)

	//we'll append our token and refresh token to this update object
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	//we'll update our updated_at field
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	//we'll use mongodb to update this object
	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	//basically we want to go to our userCollection and we want to update the user
	_, err := userCollection.UpdateOne(
		//context is basically used to update that particular user's data
		ctx,
		filter,
		bson.D{
			//$set is used in mongodb
			//here we'll set the update object
			{"$set", updateObj},
		},
		&opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		//function that takes in the token and returns interface and an error and ths returns
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), bson.ErrDecodeToNil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	//check if the token is ok or not
	if !ok {
		msg = fmt.Sprintf("The token is invalid!")
		msg = err.Error()
		return
	}

	//want to check expiration date of token
	//claims has all the information that user have
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("Token is expired")
		msg = err.Error()
		return
	}
	//we already know what are we returning claim
	return claims, msg

}
