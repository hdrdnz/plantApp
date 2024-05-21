package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Plant struct {
	Test   int           `json:"tfvcount"`
	Result []PlantResult `json:"results"`
}
type PlantResult struct {
	Name        string `json:"tfvname"`
	Image       string `json:"imageurl"`
	Description string `json:"description"`
	Uses        string `json:"uses"`
	Soil        string `json:"soil"`
	Climate     string `json:"climate"`
	Health      string `json:"health"`
}

func GetPlant(c *gin.Context) {
	plantName, _ := c.GetQuery("plant_name")
	plantGet := models.Plant{}
	plantUser := models.PlantUser{}
	db := database.GetDB()
	url := "https://www.tropicalfruitandveg.com/api/tfvjsonapi.php?tfvitem=" + plantName
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	var plant Plant
	if err := json.NewDecoder(resp.Body).Decode(&plant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	if plant.Test != 1 {
		if err := db.Where("name=?", strings.ToLower(plantName)).First(&plantGet).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
		plantUser.PlantId = plantGet.Id
		plantUser.UserId = uint(Claims["sub"].(float64))
		if err := db.Save(&plantUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
		c.JSON(http.StatusOK, plantGet)
	} else {
		plantUser.UserId = uint(Claims["sub"].(float64))
		plantUser.Image = plant.Result[0].Image
		plantUser.Name = plant.Result[0].Name
		if err := db.Save(&plantUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
		c.JSON(http.StatusOK, plant.Result)
	}

}

func GetDocs(c *gin.Context) {
	jsonVeri, _ := os.Open(("./docs/swagger.json"))
	byteValue, _ := ioutil.ReadAll(jsonVeri)
	result := make(map[string]interface{})
	_ = json.Unmarshal(byteValue, &result)
	c.JSON(http.StatusOK, result)
}
