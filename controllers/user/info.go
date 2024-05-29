package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingExample godoc
// @Summary Get Rose List
// @Schemes
// @Description Get Rose List
// @Tags rose
// @Security jwt
// @Accept json
// @Produce json
// @Success 200 {object} models.Rose
// @Success 400 string string
// @Router /mobil/rose [get]
func Rose(c *gin.Context) {
	db := database.GetDB()
	var rose []models.Rose
	if err := db.Find(&rose).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message:": "information not found.",
		})

	}
	c.JSON(http.StatusOK, rose)

}
