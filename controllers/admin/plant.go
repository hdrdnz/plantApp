package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlantDetail struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Uses        string `json:"uses"`
	Soil        string `json:"soil"`
	Climate     string `json:"climate"`
	Health      string `json:"health"`
}

// PingExample godoc
// @Summary Add plant
// @Schemes
// @Description Add plant
// @Tags admin
// @Accept json
// @Produce json
// @Param PlantPost body PlantDetail true "Post Plant"
// @Success 200 string string
// @Success 400 string string
// @Router /admin/add-plant [post]
func PlantPost(c *gin.Context) {
	var Result PlantDetail
	db := database.GetDB()
	if err := c.ShouldBindJSON(&Result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	if Result.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Description parameter cannot be empty.",
		})
		return
	}
	if Result.Uses == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Uses parameter cannot be empty.",
		})
		return
	}
	if Result.Soil == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Soil parameter cannot be empty.",
		})
		return
	}
	if Result.Climate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Climate parameter cannot be empty.",
		})
		return
	}
	if Result.Health == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Health parameter cannot be empty.",
		})
		return
	}
	plant := models.Plant{}
	plant.Name = Result.Name
	plant.Climate = Result.Climate
	plant.Soil = Result.Soil
	plant.Uses = Result.Uses
	plant.Description = Result.Description
	plant.Health = Result.Health
	if err := db.Save(&plant).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while recording.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration Successful",
	})
}

// PingExample godoc
// @Summary Get plant
// @Schemes
// @Description Get plant
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} models.Plant
// @Success 400 string string
// @Router /admin/plants [get]
func PlantGet(c *gin.Context) {
	var plant []models.Plant
	//flash := lib.NewFlash(c)
	db := database.GetDB()
	db.Find(&plant)
	c.JSON(http.StatusOK, plant)
}

// PingExample godoc
// @Summary Delete plant
// @Schemes
// @Description Delete plant
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 string string
// @Success 400 string string
// @Router /admin/plant/:plantid/delete [post]
func DeletePlant(c *gin.Context) {
	plantId := c.Param("plantid")
	if plantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "plant_id cannot be empty.",
		})
		return
	}
	plant := models.Plant{}
	db := database.GetDB()
	if err := db.Where("id=?", plantId).First(&plant).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The plant not found.",
		})
		return
	}
	if err := db.Delete(&plant).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error occurred while deleting",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}

// PingExample godoc
// @Summary Update plant
// @Schemes
// @Description Update plant
// @Tags admin
// @Accept json
// @Produce json
// @Param UpdatePlant body PlantDetail true "Update Plant"
// @Success 200 string string
// @Success 400 string string
// @Router /admin/plant/:plantid/update [post]
func UpdatePlant(c *gin.Context) {
	plantId := c.Param("plantid")
	var plantDetail PlantDetail
	db := database.GetDB()
	plant := models.Plant{}
	if err := c.ShouldBindJSON(&plantDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The plant not found.",
		})
		return
	}
	if err := db.Where("id=?", plantId).First(&plant).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The plant not found.",
		})
		return
	}
	if plantDetail.Name != "" {
		plant.Name = plantDetail.Name
	}
	if plantDetail.Climate != "" {
		plant.Climate = plantDetail.Climate
	}
	if plantDetail.Health != "" {
		plant.Health = plantDetail.Health
	}
	if plantDetail.Soil != "" {
		plant.Soil = plantDetail.Soil
	}
	if plantDetail.Description != "" {
		plant.Description = plantDetail.Description
	}
	if plantDetail.Uses != "" {
		plant.Uses = plantDetail.Uses
	}
	if err := db.Save(&plant).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Plant could not be updated",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Plant updated",
	})
}
