package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"main/models"
)

type StrukturSuhu struct {
	Id     uint    `json:"id"`
	Lokasi string  `json:"lokasi" binding:"required"`
	Suhu   float32 `json:"suhu" binding:"required"`
}

type StrukturId struct {
	Id uint `json:"id" binding:"required"`
}

func Tampil(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var modelSuhu []models.Suhu

	hasil := db.Find(&modelSuhu)

	if hasil.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  hasil.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil tampil data",
		"data":   modelSuhu,
	})
}

func Tambah(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var dataSuhu StrukturSuhu

	if err := c.ShouldBindJSON(&dataSuhu); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})

		return
	}

	modelSuhu := models.Suhu{
		Lokasi:    dataSuhu.Lokasi,
		Suhu:      dataSuhu.Suhu,
		CreatedAt: time.Now(),
	}

	hasil := db.Create(&modelSuhu)

	if hasil.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  hasil.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil tambah data",
		"data":   modelSuhu,
	})
}

func Ubah(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var dataSuhu StrukturSuhu

	if err := c.ShouldBindJSON(&dataSuhu); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})

		return
	}

	var modelSuhu models.Suhu

	hasil := db.First(&modelSuhu, dataSuhu.Id)

	if hasil.Error != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"pesan":  "Data tidak ditemukan",
		})

		return
	}

	modelSuhu.Lokasi = dataSuhu.Lokasi
	modelSuhu.Suhu = dataSuhu.Suhu

	db.Save(&modelSuhu)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil ubah data",
		"data":   modelSuhu,
	})
}

func Hapus(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var data StrukturId

	if err := c.ShouldBindJSON(&data); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})

		return
	}

	hasil := db.Delete(&models.Suhu{}, data.Id)

	if hasil.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  hasil.Error.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil hapus data",
	})
}
