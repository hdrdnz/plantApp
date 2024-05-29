package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingExample godoc
// @Summary Get User Information
// @Schemes
// @Description Get User Information by id or nickname
// @Tags user
// @Security jwt
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param nick_name query string false "User Nickname"
// @Success 200 {object} models.User
// @Success 400 string string
// @Router /mobil/user [get]
func GetUserById(c *gin.Context) {
	userId, _ := c.GetQuery("user_id")
	nickname, _ := c.GetQuery("nick_name")
	db := database.GetDB()
	user := models.User{}

	if userId == "" && nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "please enter userId or nickname",
		})
		return
	}
	if userId != "" {
		db.Where("id=?", userId).First(&user)
	}
	if nickname != "" {
		db.Where("nick_name=?", nickname).First(&user)
	}

	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "record not found",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
