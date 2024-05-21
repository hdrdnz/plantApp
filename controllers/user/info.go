package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
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

func Leaf(c *gin.Context) {
	url := "https://facts.net/leaf-facts/"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message:": "hata",
		})
		return
	}
	defer resp.Body.Close()

	// HTML içeriğini parse et
	doc, err := html.Parse(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message:": "hata2",
		})
		return
	}
	var findH2 func(*html.Node)
	i := 0
	//data := make(map[string]string)
	findH2 = func(n *html.Node) {
		ok := false
		if n.Type == html.ElementNode && n.Data == "h2" {
			ok = true
			i++
			if i != 1 {
				fmt.Println("h2 içeriği:", n.FirstChild.Data)
				fmt.Println(" n.NextSibling.Type:", n.LastChild.Data)
				if n.NextSibling.Type == html.TextNode {
					fmt.Println("data:", n.NextSibling.Data)
				}
			}

		}
		if n.Type == html.ElementNode && n.Data == "p" {
			if ok {

			}
			i++
			fmt.Println("h2 içeriği:", n.FirstChild.Data)
		}

		// if n.NextSibling != nil && n.NextSibling.Type == html.ElementNode && n.Data == "p" {
		// 	fmt.Println("Sonraki metin:", n.NextSibling.Data)
		// }
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH2(c)
		}
	}

	findH2(doc)
}
