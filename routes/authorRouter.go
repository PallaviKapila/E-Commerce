package routes

import (
	controller "github.com/PallaviKapila/E-Commerce-Project/controllers"
	"github.com/gin-gonic/gin"
)

// it accepts incoming routes and we'll give reference to gin engine
func AuthRoutes(incomingRoutes *gin.Engine) {
	//simpling defining signup and user page
	//in controller file we'll have function to control signup
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
