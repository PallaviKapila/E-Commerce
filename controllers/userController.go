package controllers

import (
	"E-Commerce Project/database"

	"github.com/go-playground/validator/v10"
	//to hash our password
)

var userCollectiob *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup()

func Login()

func GetUsers()

func GetUser()
