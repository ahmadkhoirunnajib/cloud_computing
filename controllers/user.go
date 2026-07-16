package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"main/models"
)

type StrukturUser struct {
	Id       uint   `json:"id"`
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type StrukturLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserTampil(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user []models.User

	db.Find(&user)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}

func UserTambah(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var data StrukturUser

	if err := c.ShouldBindJSON(&data); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})

		return
	}

	var sha = sha1.New()
	sha.Write([]byte(data.Password))

	passwordEncrypt := fmt.Sprintf("%x", sha.Sum(nil))

	user := models.User{
		Nama:      data.Nama,
		Username:  data.Username,
		Password:  passwordEncrypt,
		CreatedAt: time.Now(),
	}

	db.Create(&user)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}

func UserUbah(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var data StrukturUser

	c.ShouldBindJSON(&data)

	var user models.User

	db.First(&user, data.Id)

	user.Nama = data.Nama
	user.Username = data.Username

	if data.Password != "" {

		var sha = sha1.New()

		sha.Write([]byte(data.Password))

		user.Password = fmt.Sprintf("%x", sha.Sum(nil))
	}

	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}

func UserHapus(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var data StrukturUser

	c.ShouldBindJSON(&data)

	db.Delete(&models.User{}, data.Id)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil hapus user",
	})
}

func UserLogin(c *gin.Context) (interface{}, error) {

	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturLogin

	if err := c.ShouldBindJSON(&dataUser); err != nil {
		return nil, jwtV3.ErrMissingLoginValues
	}

	var sha = sha1.New()

	sha.Write([]byte(dataUser.Password))

	encrypted := fmt.Sprintf("%x", sha.Sum(nil))

	var modelUser models.User

	cekUser := db.
		Where("username = ?", dataUser.Username).
		Where("password = ?", encrypted).
		First(&modelUser)

	if cekUser.Error == nil {
		return modelUser, nil
	}

	return nil, jwtV3.ErrFailedAuthentication
}

func Profile(c *gin.Context) {

	claims := jwtV3.ExtractClaims(c)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   claims,
	})
}
