package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var Claims jwt.MapClaims

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		autheader := c.GetHeader("Authorization")
		if autheader == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please enter a valid token.",
			})
			c.Abort()
		}

		auth := strings.Split(autheader, " ")

		token, err := parseToken(auth[1])
		if err != nil {
			//token süresi kontrolü
			userToken := models.UserToken{}
			if err := db.Where("token=?", auth[1]).First(&userToken).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "user not found.",
				})
				return
			}
			if err := db.Delete(&userToken).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "user not found.",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "The token has expired.",
			})
			c.Abort()
		}
		if !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "token not found.",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//token den gelen kullanıcıyı bulma
			user := &models.User{}
			db.Where("id=?", claims["sub"]).First(user)
			if user.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "user not found.",
				})
				c.Abort()
				return
			}

			//kullanıcı-token kontrolü
			pUser := &models.UserToken{}
			db.Where("token=?", auth[1]).First(pUser)
			if pUser.UserId != user.Id {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "User-Token mismatch",
				})
				c.Abort()
				return
			}
			fmt.Println("------claims:", claims)
			Claims = claims
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "user not found.",
			})
			c.Abort()
			return
		}

	}
}
func parseToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return GetSecret(), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token.")
	}
	return token, nil

}
func GetSecret() []byte {
	file, _ := os.Open("secret.key")
	data, _ := ioutil.ReadAll(file)
	response := strings.Split(string(data), ":")
	return []byte(response[1])
}
