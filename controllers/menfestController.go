package controllers

import (
	"Bemenfest/initializers"
	"Bemenfest/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendMenfest(c *gin.Context) {
	var body struct {
		Pesan  string `json:"pesan"`
		UserID string `json:"user_id"`
	}
	// add params
	id := c.Param("id")
	// bind by params id user
	body.UserID = id
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// query user
	user := models.User{}
	result := initializers.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User tidak ditemukan"})
		return
	}
	// create menfest
	menfest := models.Menfest{
		Pesan:  body.Pesan,
		UserID: body.UserID,
	}
	result = initializers.DB.Create(&menfest)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat menfest"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Menfest berhasil dikirim"})

}

func UserMenfest(c *gin.Context) {
	// Get user and their menfests
	var user models.User
	var menfest []models.Menfest
	id := c.Param("id")
	// query user

	result := initializers.DB.Preload("Menfest").Where("id = ?", id).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User tidak ditemukan"})
		return
	}

	//query menfest
	result = initializers.DB.Where("user_id = ?", id).Find(&menfest)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menfest tidak ditemukan"})
		return
	}
	// order menfest by created at
	initializers.DB.Order("created_at desc").Find(&menfest)
	// ignore update at

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"menfest": menfest, // Access Menfest directly from the user variable
	})
}

func DeleteMenfest(c *gin.Context) {
	var menfest models.Menfest
	id := c.Param("id")
	result := initializers.DB.Where("id = ?", id).Delete(&menfest)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menfest tidak ditemukan"})
		return
	}
	// jika menfest sudah di delete
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menfest sudah dihapus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menfest berhasil dihapus"})
}
