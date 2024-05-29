package controllers

import (
	"PlantApp/database"
	"PlantApp/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string      `json:"token"`
	User    models.User `json:"user"`
	Message string      `json:"message"`
}

// PingExample godoc
// @Summary User Login
// @Schemes
// @Description Login User
// @Tags user
// @Accept json
// @Produce json
// @Param login body UserLogin true "Login user"
// @Success 200 {object} LoginResponse
// @Success 400 string string
// @Router /mobil/user/login [post]
func Login(c *gin.Context) {
	var loginUser UserLogin
	db := database.GetDB()
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	if loginUser.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or nickname cannot be  empty.",
		})
		return
	}
	user := models.User{}
	if err := db.Where("nick_name = ?", loginUser.UserName).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or password is incorrect",
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or password is incorrect",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(GetSecret())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": "something went wrong.",
		})
	}
	userToken := models.UserToken{}
	db.Where("user_id=?", user.Id).First(&userToken)
	if userToken.Id == 0 {
		userToken.UserId = user.Id
	}
	userToken.Token = tokenString
	if err := db.Save(&userToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"user":    user,
		"message": "success",
	})
}

// PingExample godoc
// @Summary User Logout
// @Schemes
// @Description Logout User
// @Tags user
// @Accept json
// @Produce json
// @Security jwt
// @Param user_id query string true "User ID"
// @Success 200 {string} Successfully logged out.
// @Success 400 string string
// @Router /mobil/user/logout [get]
func Logout(c *gin.Context) {
	userId := c.Query("user_id")
	db := database.GetDB()
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "please enter a userId.",
		})
		return
	}

	if userId != strconv.Itoa(int(Claims["sub"].(float64))) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user mismatch.",
		})
		return
	}
	userToken := models.UserToken{}
	if err := db.Where("user_id=?", userId).First(&userToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user not found.",
		})
		return
	}
	if err := db.Delete(&userToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out.",
	})
}
