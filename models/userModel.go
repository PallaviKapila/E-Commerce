// we'll call this package so that we can easily call this in our main.go file
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// this customerModel acts as a middle layer b/w our golang program and any other database
// database understands jso n and golang doesn't understand json
// we need a layer that converts things from json to golang and golang to json
// that's how we have struct here
type User struct {
	ID            primitive.ObjectID `bson:"_id`
	First_Name    *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_Name     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token"`
	User_type     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"` //enum concept
	Refresh_token *string            `json:"refresh_token`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
