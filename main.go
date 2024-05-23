package main

import (
	"PlantApp/routers"

	"github.com/gin-gonic/gin"
)

// type Test struct {
// 	Test   int         `json:"tvfcount"`
// 	Result []TestImage `json:"results"`
// }

// type TestImage struct {
// 	Name  string `json:"tfvname"`
// 	Image string `json:"imageurl"`
// }

// @title          Leaflove Mobil API
// @version         1.0
// @description     This is plantapp

// @host      https://leaflove.com.tr
func main() {
	router := gin.Default()
	router.Static("./images", "images")
	routers.Load(router)
	router.Run()
	//router.Run("localhost:9090")
	// if err := router.Run("localhost:8080"); err != nil {
	// 	fmt.Println("router hata:", err)
	// } else {
	// 	fmt.Println("başarılı")
	// }

	// var data []string
	// file, err := os.OpenFile("../imagefile.txt", os.O_WRONLY|os.O_CREATE, 0644)
	// if err != nil {
	// 	fmt.Println("okuma hatası:", err)
	// }
	// data = append(data, "Rambutan", "Rosemary", "Sapodilla", "Sesame", "Sittu", "Strawberry Guava", "Sugar Apple", "Sugarcane", "Tamarillo", "Tindora", "Tropical Almond", "Tuar", "Tulsi", "Turmeric", "Urad", "Vanilla", "Wood Apple", "Zalacca")
	// for _, name := range data {
	// 	url := "https://www.tropicalfruitandveg.com/api/tfvjsonapi.php?tfvitem=" + name
	// 	response, err := http.Get(url)
	// 	if err != nil {
	// 		fmt.Println("veri çekme hatası", err)
	// 		return
	// 	}
	// 	var data Test
	// 	defer response.Body.Close()

	// 	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
	// 		fmt.Println("API'den gelen veri JSON olarak decode edilemedi:", err)
	// 		return
	// 	}
	// 	var content string
	// 	for _, i := range data.Result {
	// 		content += i.Name + " : " + i.Image + "\n"
	// 	}
	// 	fmt.Println("content:", content)
	// 	_, err = file.Seek(0, 2) // Dosyanın sonuna git
	// 	if err != nil {
	// 		fmt.Println("Dosya konumlandırma hatası:", err)
	// 		return
	// 	}
	// 	_, err = file.Write([]byte(content))
	// 	if err != nil {
	// 		fmt.Println("yazma hatasu:", err)
	// 		return
	// 	}
	// }

	//url := "https://www.fruityvice.com/api/fruit/apple"

	//resp, err := json.Marshal(response.Body)
	// if err != nil {
	// 	fmt.Println("marshall err:", err)
	// 	return
	// }
	//data := make(map[string]interface{})

	// if err := json.Unmarshal(resp, &data); err != nil {
	// 	fmt.Println("unmarsahl:", err)
	// 	return
	// }

}
