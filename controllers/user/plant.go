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
		"Banana":   0,
		"Betelnut": 1,
		"chinar":   2,
		"Cocount":  3,
		"corn":     4,
		"grape":    5,
		"Kiwi":     6,
		"Lime":     7,
		"Mango":    8,
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

	api_key := "Df2BDaC2337z7b2Z7tr5"
	model_endpoint := "plant_detector-wknzm/1"
	//fotoğrafı açma
	f, err := os.Open("image.jpg")
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
	fruits := map[string]float64{
		"Banana":   0,
		"Betelnut": 1,
		"chinar":   2,
		"Cocount":  3,
		"corn":     4,
		"grape":    5,
		"Kiwi":     6,
		"Lime":     7,
		"Mango":    8,
		"tomato":   9,
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
	//c.JSON(http.StatusOK, string(bytes))

}

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
	c.JSON(http.StatusOK, plantUser)
}

type ImageData struct {
	ImageURL string `json:"image_url"`
}

func Images(c *gin.Context) {
	var imageData ImageData
	if err := c.ShouldBindJSON(&imageData); err != nil {
		fmt.Println("hataaa:", err)
		return
	}
	if imageData.ImageURL == "" {
		fmt.Println("boş")
		return

	}

	parts := strings.Split(imageData.ImageURL, ";base64,")
	fmt.Println("parts:", parts[0])
	if len(parts) != 2 {
		fmt.Println("Invalid Base64 data format")
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Dosya yolunu oluşturma
	deneme := strconv.FormatInt(time.Now().UnixMilli(), 10)
	filename := filepath.Base(deneme + extension)
	filePath := filepath.Join("./images", filename)

	// Dosyayı kaydetme
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "BAşarılııı",
	})

}
