// base of this project
package main

import (
	//where this routes come from
	routes "E-Commerce Project/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	//we'll have enviornment file so we'll get the enviornment
	port = os.Getenv("PORT")

	//if port is empty then atleast we need something
	if port == "" {
		port = "8000"
	}

	//gin is helping you to create router
	router := gin.New()
	router.Use(gin.Logger())

	//two types of routes authorRouter and userRouter
	//package routes
	//we'll be accessing them using routes variable and we're going to pass the gin router we created
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	//create 2 api's
	//in these router handler functions when we're using gin we don't need to pass the w and r the response and request
	//here we only need to pass gin context and we'll be able to access the respnse and request from here itself
	router.GET("/api-1", func(c *gin.Context) {
		//setting headers through gin
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})
	//port in this case is 8000 if env. board is empty otherwise it'll be whatever we define in the env. variable file
	router.Run(":" + port)
}
