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

// PingExample godoc
// @Summary User Create
// @Schemes
// @Description Create User
// @Tags user
// @Accept json
// @Produce json
// @Param register body Register true "Create user"
// @Success 200 string  "Success" "example:Success"
// @Router /mobil/user/create [post]
func UserRegister(c *gin.Context) {
	db := database.GetDB()
	var uRegister Register
	if err := c.ShouldBindJSON(&uRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	User := models.User{}
	if uRegister.NickName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username  cannot be empty.",
		})
		return
	}
	User.NickName = uRegister.NickName

	if uRegister.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email cannot be empty",
		})
		return
	}
	User.Email = uRegister.Email

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
	User.Password = string(hashedPassword)

	if err := db.Save(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message:": "Success",
	})

}
