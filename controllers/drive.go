package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"main/models"
)

func DriveUpload(c *gin.Context) {

	fileName := c.PostForm("fileName")

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "File tidak ditemukan",
			"error":  err.Error(),
		})
		return
	}

	mimeType := file.Header.Get("Content-Type")

	fileOpen, err := file.Open()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membuka file",
			"error":  err.Error(),
		})
		return
	}

	defer fileOpen.Close()

	fileData, err := io.ReadAll(fileOpen)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca file",
			"error":  err.Error(),
		})
		return
	}

	// Encode file ke Base64
	data := base64.StdEncoding.EncodeToString(fileData)

	postBody, err := json.Marshal(map[string]string{
		"fileName": fileName,
		"mimeType": mimeType,
		"data":     data,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membuat request",
			"error":  err.Error(),
		})
		return
	}

	requestBody := bytes.NewBuffer(postBody)

	// URL Google Apps Script Upload
	res, err := http.Post(
		"https://script.google.com/macros/s/AKfycbzDA7IWdpmIsSRXI7Bs5Il0EFf-aGQ0JzmOk52DKO4Q9muHJR0nhBNlkDPmiJrTye3XEA/exec",
		"application/json; charset=UTF-8",
		requestBody,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"kode":   "ERR-DRIVE",
			"pesan":  "Gagal Upload ke Google Drive",
			"error":  err.Error(),
		})
		return
	}

	defer res.Body.Close()

	hasilBody, err := io.ReadAll(res.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca response",
			"error":  err.Error(),
		})
		return
	}

	var hasilJson map[string]interface{}

	err = json.Unmarshal(hasilBody, &hasilJson)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   false,
			"pesan":    "Response bukan JSON valid",
			"response": string(hasilBody),
		})
		return
	}

	// Ambil koneksi database
	db := c.MustGet("db").(*gorm.DB)

	// Simpan metadata file ke database
	dokumenBaru := models.Dokumen{
		NamaDokumen: hasilJson["filename"].(string),
		FileId:      hasilJson["fileId"].(string),
		FileUrl:     hasilJson["fileUrl"].(string),
	}

	hasilDokumen := db.Create(&dokumenBaru)

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"pesan":     "Berhasil Upload",
		"data":      hasilJson,
		"tersimpan": hasilDokumen.RowsAffected,
	})
}

// Menampilkan semua file yang tersimpan
func DriveTampil(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var dokumen []models.Dokumen

	db.Find(&dokumen)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil Tampil",
		"data":   dokumen,
	})
}

// Download file dari Google Drive
func DriveUnduh(c *gin.Context) {

	id := c.Param("id")

	// URL Apps Script Download
	res, err := http.Get(
		"https://script.google.com/macros/s/AKfycbzDA7IWdpmIsSRXI7Bs5Il0EFf-aGQ0JzmOk52DKO4Q9muHJR0nhBNlkDPmiJrTye3XEA/exec?id=" + id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal Unduh",
		})
		return
	}

	defer res.Body.Close()

	hasilBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca response",
		})
		return
	}

	var hasilJson map[string]interface{}

	err = json.Unmarshal(hasilBody, &hasilJson)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Response tidak valid",
		})
		return
	}

	fileBase64 := hasilJson["file"].(string)
	mimeType := hasilJson["mimeType"].(string)

	fileData, err := base64.StdEncoding.DecodeString(fileBase64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal decode file",
		})
		return
	}

	c.Writer.Header().Set("Content-Type", mimeType)
	c.Writer.Write(fileData)
}
