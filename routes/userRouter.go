package routes

import (
	controller "E-Commerce Project/controllers"
	"E-Commerce Project/middleware"

	"githib.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	//here we're using middleware to ensure these both are protected routes
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
