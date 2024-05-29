package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebPlant struct {
	Name        string
	Image       string
	Description string
}

// PingExample godoc
// @Summary Web Plants
// @Schemes
// @Description Web Plants
// @Tags web
// @Accept json
// @Produce json
// @Success 200 {array} WebPlant
// @Success 400 string string
// @Router /web/plants [get]
func WebPlants(c *gin.Context) {
	db := database.GetDB()
	var rose []models.Rose
	var wPlants []WebPlant
	if err := db.Find(&rose).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message:": "information not found.",
		})

	}
	for _, item := range rose {
		var plant WebPlant
		plant.Name = item.Name
		plant.Description = item.Description
		plant.Image = item.Image
		wPlants = append(wPlants, plant)
	}
	var generalPlants []models.GeneralPlants
	if err := db.Find(&generalPlants).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message:": "information not found.",
		})

	}
	for _, item := range generalPlants {
		var plant WebPlant
		plant.Name = item.Name
		plant.Description = item.Description
		plant.Image = item.Image
		wPlants = append(wPlants, plant)
	}
	c.JSON(http.StatusOK, wPlants)
}
