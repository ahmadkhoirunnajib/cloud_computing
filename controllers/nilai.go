package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrukturNilai struct {
	Id         uint    `json:"id"`
	NamaMhs    string  `json:"nama_mhs"`
	NilaiTugas float64 `json:"nilai_tugas"`
	NilaiUts   float64 `json:"nilai_uts"`
	NilaiUas   float64 `json:"nilai_uas"`
}

func hitungNilai(tugas, uts, uas float64) (float64, string, string) {

	nilaiAkhir := (tugas * 0.30) +
		(uts * 0.30) +
		(uas * 0.40)

	var huruf string
	var status string

	if nilaiAkhir >= 80 {

		huruf = "A"
		status = "LULUS"

	} else if nilaiAkhir >= 65 {

		huruf = "B"
		status = "LULUS"

	} else if nilaiAkhir >= 50 {

		huruf = "C"
		status = "MENGULANG"

	} else {

		huruf = "D"
		status = "MENGULANG"
	}

	return nilaiAkhir, huruf, status
}

func NilaiTambah(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var data StrukturNilai

	if err := c.ShouldBindJSON(&data); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})

		return
	}

	angka, huruf, status := hitungNilai(
		data.NilaiTugas,
		data.NilaiUts,
		data.NilaiUas,
	)

	nilai := models.Nilai{
		NamaMhs:    data.NamaMhs,
		NilaiTugas: data.NilaiTugas,
		NilaiUts:   data.NilaiUts,
		NilaiUas:   data.NilaiUas,
		NilaiAngka: angka,
		NilaiHuruf: huruf,
		Status:     status,
	}

	db.Create(&nilai)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   nilai,
	})
}

func NilaiTampil(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var nilai []models.Nilai

	db.Find(&nilai)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   nilai,
	})
}

func NilaiUbah(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var data StrukturNilai

	if err := c.ShouldBindJSON(&data); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})

		return
	}

	var nilai models.Nilai

	db.First(&nilai, data.Id)

	angka, huruf, status := hitungNilai(
		data.NilaiTugas,
		data.NilaiUts,
		data.NilaiUas,
	)

	nilai.NamaMhs = data.NamaMhs
	nilai.NilaiTugas = data.NilaiTugas
	nilai.NilaiUts = data.NilaiUts
	nilai.NilaiUas = data.NilaiUas

	nilai.NilaiAngka = angka
	nilai.NilaiHuruf = huruf
	nilai.Status = status

	db.Save(&nilai)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   nilai,
	})
}

func NilaiHapus(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var data struct {
		Id uint `json:"id"`
	}

	c.ShouldBindJSON(&data)

	db.Delete(&models.Nilai{}, data.Id)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil hapus nilai",
	})
}
