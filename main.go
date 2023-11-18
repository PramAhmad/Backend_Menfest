package main

import (
	"Bemenfest/controllers"
	"Bemenfest/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()
	// allow cross for everything
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	// menfest
	r.POST("/menfest/:id", controllers.SendMenfest)

	// Auth
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// profile
	r.GET("/profile/:id", controllers.UserProfile)
	// menfest per user
	r.GET("/menfest/:id", controllers.UserMenfest)
	r.DELETE("/menfest/:id", controllers.DeleteMenfest)
	r.Run()
}
