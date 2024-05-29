package controllers

import (
	controllers "PlantApp/controllers/user"
	"PlantApp/database"
	"PlantApp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LRespAdmin struct {
	Token   string      `json:"token"`
	Admin   models.User `json:"user"`
	Message string      `json:"message"`
}

// PingExample godoc
// @Summary Admin Login
// @Schemes
// @Description Login Admin
// @Tags admin
// @Accept json
// @Produce json
// @Param login body controllers.UserLogin true "Login Admin"
// @Success 200 {object} LRespAdmin
// @Success 400 string string
// @Router /admin/login [post]
func AdminLogin(c *gin.Context) {
	var loginUser controllers.UserLogin
	db := database.GetDB()
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong.",
		})
		return
	}
	if loginUser.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or nickname cannot be empty.",
		})
		return
	}
	user := models.Admin{}
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
	tokenString, err := token.SignedString(controllers.GetSecret())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": "something went wrong.",
		})
	}
	userToken := models.AdminToken{}
	db.Where("admin_id=?", user.Id).First(&userToken)
	if userToken.Id == 0 {
		userToken.AdminId = user.Id
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
		"admin":   user,
		"message": "success",
	})
}
