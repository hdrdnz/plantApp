package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	NickName string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminRegister(c *gin.Context) {
	db := database.GetDB()
	var uRegister Register
	if err := c.ShouldBindJSON(&uRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	AUser := models.Admin{}
	if uRegister.NickName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username  cannot be empty.",
		})
		return
	}
	AUser.NickName = uRegister.NickName

	if uRegister.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email cannot be empty",
		})
		return
	}
	AUser.Email = uRegister.Email

	if uRegister.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password cannot be empty.",
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(uRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	AUser.Password = string(hashedPassword)

	if err := db.Save(&AUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message:": "Success",
	})
}
