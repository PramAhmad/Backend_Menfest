package main

import (
	"Bemenfest/initializers"
	"Bemenfest/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Menfest{})
	initializers.DB.AutoMigrate(&models.User{})

}
