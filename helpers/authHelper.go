package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// helper function we're passing role here i.e user type which is admin or user
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access this resource")
		return err
	}

	return err
}

// then it returns error
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	//bcoz we want user to only access his own user data, user can't access the user data of any other user
	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access this resource!")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
