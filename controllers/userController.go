package controllers

import (
	"Bemenfest/initializers"
	"Bemenfest/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	// select user
	result := initializers.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})

}
