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
		fmt.Println("hataaa:", err)
		return
	}
	if imageData.ImageURL == "" {
		fmt.Println("boş")
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
	f, err := os.Open("image.jpg")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	//bitki ismi için istek atma
	data := base64.StdEncoding.EncodeToString(content)
	fmt.Println("data :", data)
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

func PostPlant(c *gin.Context) {
	plant := []models.Plant{
		{
			Name:        "Chinar",
			Image:       "https://ideacdn.net/shop/hr/14/myassets/products/791/cinar-agaci-fidani-001.jpg",
			Description: "The chinar is a broadleaf, long-lived, and fast-growing tree. It is often used in parks and gardens for providing shade.",
			Uses:        "Sycamore trees are primarily used for shade. Additionally, their wood is valuable for furniture and building materials.",
			Soil:        "Sycamore trees prefer well-drained, fertile, and moist soils. The soil pH should be between 6.0 and 7.5.",
			Climate:     "They thrive in temperate climates. They prefer humid environments and grow better near water bodies.",
			Health:      "Sycamore trees are generally hardy but are susceptible to diseases such as anthracnose, leaf spot, and powdery mildew.",
		},
		{
			Name:        "Corn",
			Image:       "https://cdn.britannica.com/36/167236-050-BF90337E/Ears-corn.jpg",
			Description: "Corn is a cereal plant that is one of the most important and widely grown crops. It has large, elongated ears and kernels.",
			Uses:        "Corn is used for human consumption, as livestock feed, and in industrial products like ethanol, corn syrup, and bioplastics.",
			Soil:        "Corn prefers well-drained, fertile soils with a pH between 5.8 and 7.0. It requires high levels of nitrogen.",
			Climate:     "Corn thrives in warm climates with plenty of sunlight. It requires moderate to high rainfall for optimal growth.",
			Health:      "Corn plants are susceptible to pests and diseases such as corn borers, rootworms, and fungal diseases like rust and blight.",
		},
		{
			Name:        "Grape",
			Image:       "https://esular.com/wp-content/uploads/2022/04/green-grapes.jpg",
			Description: "Grapes are a type of fruit that grow in clusters and can be eaten fresh or used to make wine, juice, and raisins.",
			Uses:        "Grapes are consumed fresh, dried (as raisins), or processed into wine, juice, and various culinary products.",
			Soil:        "Grapes prefer well-drained, loamy soils with a pH between 5.5 and 6.5. They require good air circulation to prevent diseases.",
			Climate:     "Grapes thrive in temperate climates with warm, dry summers and mild winters. They require plenty of sunlight for optimal growth.",
			Health:      "Grape plants are susceptible to pests and diseases such as powdery mildew, downy mildew, and phylloxera.",
		},
		{
			Name:        "Tomato",
			Image:       "https://i.lezzet.com.tr/images-xxlarge-secondary/domates-meyve-mi-sebze-mi-68f04593-f07c-43de-afe2-bc9948bcfbdd.jpg",
			Description: "Tomatoes (Solanum lycopersicum) are a widely cultivated fruit-bearing plant, known for their bright red, juicy, and flavorful fruits.",
			Uses:        "Tomatoes are used in a variety of culinary applications including salads, sauces, soups, and as a base for many dishes.",
			Soil:        "Tomatoes prefer well-drained, fertile soils rich in organic matter, with a pH between 6.0 and 6.8.",
			Climate:     "Tomatoes thrive in warm, sunny climates with temperatures between 70°F and 85°F (21°C to 29°C). They require consistent moisture but do not tolerate waterlogging.",
			Health:      "Tomato plants are susceptible to pests and diseases such as tomato hornworms, aphids, blight, and blossom end rot.",
		},
	}

	db := database.GetDB()
	if err := db.Save(&plant).Error; err != nil {
		fmt.Println("hataaaa:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "başarılıı",
	})

}
func GeneralPlant(c *gin.Context) {
	generalPlants := []models.GeneralPlants{
		{
			Name:        "Yacon",
			Image:       "https://iatkv.tmgrup.com.tr/bc679e/0/0/0/0/640/440?u=https://itkv.tmgrup.com.tr/2018/12/24/1545642267080.jpeg&mw=616",
			Description: "Yacon is a sweet and juicy root vegetable. It has a low glycemic index and prebiotic properties.",
		},
		{
			Name:        "Rosehip",
			Image:       "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT6TkwcP5IkqxNiZMFoBEL5Sc7ZA1psTQZpRuSdzf86rg&s",
			Description: "Rosehip is the fruit of plants belonging to the rose family. It contains high levels of vitamin C and is commonly used in making tea or jam.",
		},
		{
			Name:        "Chayote",
			Image:       "https://images.ctfassets.net/ruek9xr8ihvu/55E60jljWy3G3KUr74WDO4/9cd0e0283fe61a29a91b26545820b18c/Chayote-for-Babies-scaled.jpg?w=2000&q=80&fm=webp",
			Description: "Chayote is a vegetable belonging to the squash family. It is consumed cooked as a side dish or in salads.",
		},
		{
			Name:        "Jicama",
			Image:       "https://cdn-prod.medicalnewstoday.com/content/images/articles/324/324241/jicama-on-a-table.jpg",
			Description: "Jicama is a root vegetable native to Mexico. It has a crunchy texture and a slightly sweet taste. It is consumed in salads or eaten raw as a snack.",
		},
		{
			Name:        "Rambutan",
			Image:       "https://tropikalmeyveler.com/wp-content/uploads/2019/03/rambutan.jpeg",
			Description: "Rambutan is a tropical fruit with a hairy outer shell. The flesh inside is sweet and juicy, and it is often consumed fresh.",
		},
		{
			Name:        "Nopal Cactus",
			Image:       "https://media.post.rvohealth.io/wp-content/uploads/2020/08/732x549_THUMBNAIL_Nopal_Cactus.jpg",
			Description: "Nopal cactus is a type of cactus that grows in Mexico and Central America. Its leaves are edible and are often used in salads or cooked dishes.",
		},
		{
			Name:        "Biriba",
			Image:       "https://static.wixstatic.com/media/cd8fa9_34631bf9e086414590937228c9963553~mv2.jpeg/v1/fill/w_640,h_368,al_c,q_80,enc_auto/cd8fa9_34631bf9e086414590937228c9963553~mv2.jpeg",
			Description: "Biriba is a fruit native to Brazil, it has a sweet and tangy flavor. It contains vitamin C, vitamin A, and iron.",
		},
	}
	db := database.GetDB()
	if err := db.Save(&generalPlants).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "başarılı kayıt.",
	})

}

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
	err := uploadImage(imageData.ImageURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
}

func GetFavorites(c *gin.Context) {
	userId, _ := c.GetQuery("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id cannot be empty.",
		})
		return
	}
	db := database.GetDB()
	favorite := models.Favorites{}
	if err := db.Where("user_id=?", userId).Preload("Rose").Preload("GeneralPlants").Preload("user").First(&favorite).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "favorite list not found.",
		})
		return
	}
	c.JSON(http.StatusOK, favorite)
}

func uploadImage(imageUrl string) error {
	parts := strings.Split(imageUrl, ";base64,")
	fmt.Println("parts:", parts[0])
	if len(parts) != 2 {
		return fmt.Errorf("Invalid Base64 data format")
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
		return err
	}

	// Dosya yolunu oluşturma
	deneme := strconv.FormatInt(time.Now().UnixMilli(), 10)
	filename := filepath.Base(deneme + extension)
	filePath := filepath.Join("./images", filename)

	// Dosyayı kaydetme
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("Failed to save image")
	}
	return nil
}
