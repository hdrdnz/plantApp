package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	image := c.PostForm("image_url")
	api_key := "xtiSqLrecnZPeGCZtRS6"
	model_endpoint := "plant_detector-ggbfz/1"
	f, err := os.Open(image)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	// Encode as base64.
	data := base64.StdEncoding.EncodeToString(content)
	uploadURL := "https://detect.roboflow.com/" + model_endpoint + "?api_key=" + api_key + "&name=bitki2.jpeg"
	req, _ := http.NewRequest("POST", uploadURL, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	if err := json.Unmarshal([]byte(bytes), &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}

	// 0 muz +
	// 1 betel betelnut
	// 2chinar
	// 3coconut +
	// 4corn
	// 5grape
	// 6kiwi +
	// 7lime +
	// 8mango +
	// 9tomato
	fruits := map[string]int{
		"banana":   0,
		"betelnut": 1,
		"chinar":   2,
		"cocount":  3,
		"corn":     4,
		"grape":    5,
		"kiwi":     6,
		"lime":     7,
		"tomato":   9,
	}
	var plantName string
	for item, value := range fruits {
		if response["class_id"] == value {
			plantName = item
		}
	}
	//plantName, _ := c.GetQuery("plant_name")
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

func GetImage(c *gin.Context) {
	//image := c.PostForm("image")
	api_key := "xtiSqLrecnZPeGCZtRS6"          // Your API Key
	model_endpoint := "plant_detector-ggbfz/1" // Set model endpoint
	// Open file on disk.    // gorsel ismi buraya girilecek
	// Read entire JPG into byte slice.

	//f, err := os.Open(image)
	f, err := os.Open("bitki2.jpeg")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	// Encode as base64.
	data := base64.StdEncoding.EncodeToString(content)
	uploadURL := "https://detect.roboflow.com/" + model_endpoint + "?api_key=" + api_key + "&name=bitki.jpeg"
	req, _ := http.NewRequest("POST", uploadURL, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
	c.JSON(http.StatusOK, string(bytes))

}
