package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	UploadImage string
}

func GetPlant(c *gin.Context) {
	//fotoğraf alma
	imageURL := c.Query("image_url")
	//resim indirme
	i_response, err := http.Get(imageURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to download image",
		})
	}
	defer i_response.Body.Close()
	//hedef dosya
	file, err := os.Create("image.jpg")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create file",
		})
	}
	defer file.Close()

	// Response body'den dosyaya veri kopyala
	_, err = io.Copy(file, i_response.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to save image",
		})
	}
	api_key := "xtiSqLrecnZPeGCZtRS6"
	model_endpoint := "plant_detector-ggbfz/1"
	//fotoğr
	f, err := os.Open("image.jpg")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	// Encode as base64.
	data := base64.StdEncoding.EncodeToString(content)
	uploadURL := "https://detect.roboflow.com/" + model_endpoint + "?api_key=" + api_key + "&name=image.jpg"
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
	// 1 betel betelnut+
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
		"Betelnut": 1,
		"chinar":   2,
		"coconut":  3,
		"corn":     4,
		"grape":    5,
		"Lime":     6,
		"Mango":    7,
		"tomato":   8,
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
	var imageData ImageData
	if err := c.ShouldBindJSON(&imageData); err != nil {
		return
	}
	if imageData.ImageURL == "" {
		return

	}
	//fotoğraf alma
	imageURL := imageData.ImageURL
	//resim indirme
	i_response, err := http.Get(imageURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to download image",
		})
	}
	defer i_response.Body.Close()
	//hedef dosya
	file, err := os.Create("image.jpg")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create file",
		})
	}
	defer file.Close()

	// Response body'den dosyaya veri kopyala
	_, err = io.Copy(file, i_response.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to save image",
		})
	}

	api_key := "Df2BDaC2337z7b2Z7tr5"
	model_endpoint := "plant_detector-wknzm/1"
	//fotoğrafı açma
	f, err := os.Open("filenam")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	//bitki ismi için istek atma
	data := base64.StdEncoding.EncodeToString(content)
	uploadURL := "https://detect.roboflow.com/" + model_endpoint + "?api_key=" + api_key + "&name=image.jpg"
	req, _ := http.NewRequest("POST", uploadURL, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	//sonuç kısmının alınması
	if err := json.Unmarshal([]byte(bytes), &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	predictions, ok := response["predictions"].([]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	var ImageID float64
	for _, prediction := range predictions {
		predMap, ok := prediction.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
		ImageID = predMap["class_id"].(float64)
	}
	// 0 muz +
	// 1 betel betelnut+
	// 2chinar
	// 3coconut +
	// 4corn
	// 5grape
	// 6lime +
	// 7mango +
	// 8tomato
	//c.JSON(http.StatusOK, string(bytes))
	fruits := map[string]float64{
		"Banana":   0,
		"Betelnut": 1,
		"chinar":   2,
		"Coconut":  3,
		"corn":     4,
		"grape":    5,
		"Lime":     6,
		"Mango":    7,
		"tomato":   8,
	}
	//bitki isminin alınması
	var plantName string
	for item, value := range fruits {
		if ImageID == value {
			plantName = item
		}
	}

	plantGet := models.Plant{}
	plantUser := models.PlantUser{}
	db := database.GetDB()
	url := "https://www.tropicalfruitandveg.com/api/tfvjsonapi.php?tfvitem=" + plantName
	respx, err := http.Get(url)
	defer respx.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	var plant Plant
	if err := json.NewDecoder(respx.Body).Decode(&plant); err != nil {
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
		plantUser.Image = imageURL
		plantUser.Name = plantName
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
		plantUser.Image = imageURL
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

type MyPlant struct {
	PlantUserId uint
	Name        string
	Image       string
}

// PingExample godoc
// @Summary Get User Plants
// @Schemes
// @Description Get User Plant by user_id for my plants section
// @Tags plant
// @Security jwt
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Success 200 {array} MyPlant
// @Success 400 string string
// @Router /mobil/plants [get]
func MyPlants(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id cannot be empty.",
		})
	}
	if userId != strconv.Itoa(int(Claims["sub"].(float64))) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user not found.",
		})
	}
	var plantUser []models.PlantUser
	db := database.DB
	if err := db.Where("user_id=?", userId).Preload("User").Preload("Plant").Find(&plantUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
	}
	if len(plantUser) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No record found.",
		})
	}
	var myPlant []MyPlant
	for _, item := range plantUser {
		myPlant = append(myPlant, MyPlant{
			PlantUserId: item.Id,
			Name:        item.Name,
			Image:       item.Image,
		})
	}
	c.JSON(http.StatusOK, myPlant)
}

type ImageData struct {
	ImageURL string `json:"image_url"`
}

// PingExample godoc
// @Summary Get General Plants
// @Schemes
// @Description Get General Plants
// @Tags general_plants
// @Security jwt
// @Accept json
// @Produce json
// @Success 200 {array} models.GeneralPlants
// @Success 400 string string
// @Router /mobil/general-plants [get]
func GetFruits(c *gin.Context) {
	db := database.GetDB()
	var fruits []models.GeneralPlants
	if err := db.Find(&fruits).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fruits not found.",
		})
	}
	c.JSON(http.StatusOK, fruits)
}

// PingExample godoc
// @Summary Upload Plant Image
// @Schemes
// @Description Upload Plant Image by base64 parameter
// @Tags plant
// @Security jwt
// @Accept json
// @Produce json
// @Param imageData body ImageData true "Upload Plant Image"
// @Success 200 {object} models.Plant
// @Success 400 string string
// @Router /mobil/plant-upload [post]
func GetInformation(c *gin.Context) {
	var imageData ImageData
	if err := c.ShouldBindJSON(&imageData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
		})
		return
	}
	if imageData.ImageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "image_url cannot be empty",
		})
		return
	}
	//resmi kaydetme
	filename, err := uploadImage(imageData.ImageURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	api_key := "Df2BDaC2337z7b2Z7tr5"
	model_endpoint := "plant_detector-wknzm/1"
	//fotoğrafı açma
	f, err := os.Open("./images/" + filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	//bitki ismi için istek atma
	data := base64.StdEncoding.EncodeToString(content)
	uploadURL := "https://detect.roboflow.com/" + model_endpoint + "?api_key=" + api_key + "&name=" + filename
	req, _ := http.NewRequest("POST", uploadURL, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	//sonuç kısmının alınması
	if err := json.Unmarshal([]byte(bytes), &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	predictions, ok := response["predictions"].([]interface{})
	if !ok || len(predictions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "This plant is not registered in the system",
		})
		return
	}
	var ImageID float64
	for _, prediction := range predictions {
		predMap, ok := prediction.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
		ImageID = predMap["class_id"].(float64)
	}
	fruits := map[string]float64{
		"Banana":   0,
		"Betelnut": 1,
		"chinar":   2,
		"Coconut":  3,
		"corn":     4,
		"grape":    5,
		"Lime":     6,
		"Mango":    7,
		"tomato":   8,
	}
	//bitki isminin alınması
	var plantName string
	for item, value := range fruits {
		if ImageID == value {
			plantName = item
		}
	}

	plantGet := models.Plant{}
	plantUser := models.PlantUser{}
	db := database.GetDB()
	url := "https://www.tropicalfruitandveg.com/api/tfvjsonapi.php?tfvitem=" + plantName
	respx, err := http.Get(url)
	defer respx.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	var plant Plant
	if err := json.NewDecoder(respx.Body).Decode(&plant); err != nil {
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
		plantGet.Image = "https://leaflove.com.tr/images/" + filename
		// GetPlant.Name = plantGet.Name
		// GetPlant.Description = plantGet.Description
		// GetPlant.Health = plantGet.Health
		// GetPlant.Uses = plantGet.Uses
		// GetPlant.Climate = plantGet.Climate
		// GetPlant.Soil = plantGet.Soil
		// GetPlant.Image = plantGet.Image
		// GetPlant.UploadImage = "https://leaflove.com.tr/images/" + filename

		plantUser.PlantId = plantGet.Id
		plantUser.Image = "https://leaflove.com.tr/images/" + filename
		plantUser.Name = plantName
		plantUser.UserId = uint(Claims["sub"].(float64))
		if err := db.Save(&plantUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
	} else {

		plantGet.Name = plant.Result[0].Name
		plantGet.Description = plant.Result[0].Description
		plantGet.Health = plant.Result[0].Health
		plantGet.Uses = plant.Result[0].Uses
		plantGet.Climate = plant.Result[0].Climate
		plantGet.Soil = plant.Result[0].Soil
		plantGet.Image = "https://leaflove.com.tr/images/" + filename

		//plantuser kısmına kaydetme
		plantUser.UserId = uint(Claims["sub"].(float64))
		plantUser.Image = "https://leaflove.com.tr/images/" + filename
		plantUser.Name = plant.Result[0].Name
		if err := db.Save(&plantUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
	}
	c.JSON(http.StatusOK, plantGet)

}

// PingExample godoc
// @Summary Get PLant Detail
// @Schemes
// @Description Get Plant DEtail by plant_user_id
// @Tags plant
// @Security jwt
// @Accept json
// @Produce json
// @Param plant_user_id query string false "Plant user id"
// @Success 200 {object} PlantResult
// @Success 400 string string
// @Router /mobil/plant-detail [get]
func GetDetail(c *gin.Context) {
	PUserId := c.Query("plant_user_id")
	plantUser := models.PlantUser{}
	db := database.GetDB()
	if err := db.Where("id=?", PUserId).First(&plantUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	url := "https://www.tropicalfruitandveg.com/api/tfvjsonapi.php?tfvitem=" + plantUser.Name
	respx, err := http.Get(url)
	defer respx.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	var plant Plant
	if err := json.NewDecoder(respx.Body).Decode(&plant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	gPlant := models.Plant{}
	if plant.Test != 1 {
		if err := db.Where("id=?", plantUser.PlantId).First(&gPlant).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "record not found.",
			})
			return
		}
		gPlant.Image = plantUser.Image
	} else {
		gPlant.Name = plant.Result[0].Name
		gPlant.Description = plant.Result[0].Description
		gPlant.Health = plant.Result[0].Health
		gPlant.Uses = plant.Result[0].Uses
		gPlant.Climate = plant.Result[0].Climate
		gPlant.Soil = plant.Result[0].Soil
		gPlant.Image = plantUser.Image
	}
	c.JSON(http.StatusOK, gPlant)
}

type MyFav struct {
	PlantId     uint
	Name        string
	Image       string
	Description string
}

// PingExample godoc
// @Summary Get Favorites
// @Schemes
// @Description Get Favorites by user_id
// @Tags favorite
// @Security jwt
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Success 200 {array} MyFav
// @Success 400 string string
// @Router /mobil/favorites [get]
func GetFavorites(c *gin.Context) {
	var myFav []MyFav
	userId, _ := c.GetQuery("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id cannot be empty.",
		})
		return
	}
	db := database.GetDB()
	var favorite []models.Favorite
	if err := db.Where("user_id=?", userId).Preload("Rose").Find(&favorite).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "favorite list not found.",
		})
		return
	}
	for _, item := range favorite {
		if item.RoseId != 0 {
			rose := models.Rose{}
			db.Where("id=?", item.RoseId).First(&rose)
			myFav = append(myFav, MyFav{
				PlantId:     rose.Id,
				Name:        rose.Name,
				Image:       rose.Image,
				Description: rose.Description,
			})
		}
		if item.GeneralPlantsId != 0 {
			gPlant := models.GeneralPlants{}
			db.Where("id=?", item.GeneralPlantsId).First(&gPlant)
			myFav = append(myFav, MyFav{
				PlantId:     gPlant.Id,
				Name:        gPlant.Name,
				Image:       gPlant.Image,
				Description: gPlant.Description,
			})
		}
	}
	c.JSON(http.StatusOK, myFav)
}

type Fav struct {
	RoseId          uint `json:"rose_id"`
	GeneralPlantsId uint `json:"general_plants_id"`
}

// PingExample godoc
// @Summary Add Favorite
// @Schemes
// @Description Add Favorite, you need send one parameter in the parameters.
// @Tags favorite
// @Security jwt
// @Accept json
// @Produce json
// @Param fav body Fav true "Add Favorite"
// @Success 200 {string} string
// @Success 400 string string
// @Router /mobil/add-favorite [post]
func AddFavorite(c *gin.Context) {
	var Favo Fav
	db := database.GetDB()
	myFavor := models.Favorite{}
	userId := Claims["sub"].(float64)
	if err := c.ShouldBindJSON(&Favo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	var count int64
	if Favo.GeneralPlantsId != 0 {
		db.Table("favorites").Where("general_plants_id=?", Favo.GeneralPlantsId).Where("user_id=?", userId).Count(&count)
		if count != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "The plant is already in your favorite list.",
			})
			return
		}
		myFavor.GeneralPlantsId = Favo.GeneralPlantsId
	} else if Favo.RoseId != 0 {
		db.Table("favorites").Where("rose_id=?", Favo.RoseId).Where("user_id=?", userId).Count(&count)
		if count != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "The plant is already in your favorite list.",
			})
			return
		}
		myFavor.RoseId = Favo.RoseId
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please choose plant.",
		})
		return
	}
	myFavor.UserId = uint(Claims["sub"].(float64))
	if err := db.Save(&myFavor).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "add to favorite successfully.",
	})

}

type DFav struct {
	PlantName string `json:"plant_name"`
}

// PingExample godoc
// @Summary Delete Favorite
// @Schemes
// @Description Delete Favorite
// @Tags favorite
// @Security jwt
// @Accept json
// @Produce json
// @Param plant_name query string false "Plant Name"
// @Success 200 {string} string
// @Router /mobil/delete-favorite [post]
func DelFavor(c *gin.Context) {
	var dfav DFav
	if err := c.ShouldBindJSON(&dfav); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	db := database.GetDB()
	userId := Claims["sub"].(float64)
	rose := models.Rose{}
	gPlant := models.GeneralPlants{}
	var count int64
	db.Table("roses").Where("name=?", dfav.PlantName).Count(&count)
	if count != 0 {
		db.Where("name=?", dfav.PlantName).First(&rose)
	}
	fav := models.Favorite{}
	if rose.Id != 0 {
		if err := db.Where("rose_id=?", rose.Id).Where("user_id=?", userId).First(&fav).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}

	} else {
		if err := db.Where("name=?", dfav.PlantName).First(&gPlant).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
		if err := db.Where("general_plants_id=?", gPlant.Id).Where("user_id=?", userId).First(&fav).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong.",
			})
			return
		}
	}
	if fav.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "record not found.",
		})
		return
	}
	if err := db.Delete(&fav).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not be removed from favorites.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "record deleted from favorites",
	})

}

func uploadImage(imageUrl string) (string, error) {
	parts := strings.Split(imageUrl, ";base64,")
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid Base64 data format")
	}

	// MIME türünden uzantıyı belirle
	var extension string
	switch {
	case strings.HasPrefix(parts[0], "data:image/png"):
		extension = ".png"
	case strings.HasPrefix(parts[0], "data:image/jpeg"):
		extension = ".jpg"
	default:
		extension = ".bin" // Varsayılan olarak bin dosyası olarak kabul edin
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// Dosya yolunu oluşturma
	deneme := strconv.FormatInt(time.Now().UnixMilli(), 10)
	filename := filepath.Base(deneme + extension)
	filePath := filepath.Join("./images", filename)

	// Dosyayı kaydetme
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("Failed to save image")
	}
	return filename, nil
}
