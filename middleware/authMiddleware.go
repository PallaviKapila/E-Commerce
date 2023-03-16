package middleware

import (
	"fmt"
	"net/http"
	helper "github.com/PallaviKapila/E-Commerce-Project/helpers"
	"github.com/gin-gonic/gin"
)

//this function is basically taking a token from header
//whenever we call let say user routes these routes are predicted and private routes
//before these routes are called we need to authenticate so that's why we are using middleware
//which ensures all the requests are authenticated only then we can call these api's 
//signup should be public so we nned not authenticate that
func Authenticate() gin.HandlerFunc{
	return func( c *gin.Context){
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("No Authorization Header provided")})
			//to abort operations
			c.Abort()
			return
		}
		//to validate token
		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":})
			c.Abort()
			return
		}
		//we'll set the email with claims.Email, we're setting our context with these details
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("Last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Set("user_type",claims.User_type)
		c.Next()
	}
}