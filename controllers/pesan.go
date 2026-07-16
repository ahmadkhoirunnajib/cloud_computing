package controllers

import (
	"net/http"

	"main/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ======================================================
// Mengambil koneksi database dari middleware
// ======================================================

func getDB(c *gin.Context) *gorm.DB {

	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		panic("Database tidak ditemukan")
	}

	return db
}

// ======================================================
// Digunakan oleh WhatsApp Bot
// Mencari balasan berdasarkan kode
// ======================================================

func CariPesan(db *gorm.DB, kode string) (string, bool) {

	var pesan models.Pesan

	err := db.
		Where("LOWER(kode)=LOWER(?)", kode).
		First(&pesan).Error

	if err != nil {
		return "", false
	}

	return pesan.Balasan, true
}

// ======================================================
// GET /backend/pesan
// ======================================================

func PesanTampil(c *gin.Context) {

	db := getDB(c)

	var data []models.Pesan

	if err := db.Find(&data).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data",
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil",
		"data":    data,
	})
}

// ======================================================
// POST /backend/pesan
// ======================================================

func PesanTambah(c *gin.Context) {

	db := getDB(c)

	var input models.Pesan

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	var cek models.Pesan

	if err := db.Where("kode=?", input.Kode).First(&cek).Error; err == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Kode sudah digunakan",
		})

		return
	}

	if err := db.Create(&input).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil ditambahkan",
		"data":    input,
	})
}

// ======================================================
// PUT /backend/pesan
// ======================================================

func PesanUbah(c *gin.Context) {

	db := getDB(c)

	var input models.Pesan

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	var pesan models.Pesan

	if err := db.Where("kode=?", input.Kode).First(&pesan).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Kode tidak ditemukan",
		})

		return
	}

	pesan.Balasan = input.Balasan

	if err := db.Save(&pesan).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diupdate",
		"data":    pesan,
	})
}

// ======================================================
// DELETE /backend/pesan?kode=menu
// ======================================================

func PesanHapus(c *gin.Context) {

	db := getDB(c)

	kode := c.Query("kode")

	if kode == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter kode wajib diisi",
		})

		return
	}

	var pesan models.Pesan

	if err := db.Where("kode=?", kode).First(&pesan).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data tidak ditemukan",
		})

		return
	}

	if err := db.Delete(&pesan).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil dihapus",
	})
}
