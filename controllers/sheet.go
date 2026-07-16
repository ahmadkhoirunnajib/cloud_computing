package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MahasiswaSheet struct {
	ID     string  `json:"id"`
	Nama   string  `json:"nama"`
	Prodi  string  `json:"prodi"`
	Nilai  float64 `json:"nilai"`
	Action string  `json:"action,omitempty"`
}

func SheetTampil(c *gin.Context) {

	res, err := http.Get(
		"https://script.google.com/macros/s/AKfycbxdeA7wPYLvqlJxwOLoIXhw0tQ1FdS2KAzrt0hVxMW-Ppq_VMZIrmUSOZh9gmzG8DVCtw/exec",
	)

	if err != nil {

		c.JSON(500, gin.H{
			"status": false,
		})

		return
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var data interface{}

	json.Unmarshal(body, &data)

	c.JSON(200, gin.H{
		"status": true,
		"data":   data,
	})
}

func SheetTambah(c *gin.Context) {
	var mahasiswa MahasiswaSheet

	if err := c.ShouldBindJSON(&mahasiswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})
		return
	}

	postBody, _ := json.Marshal(mahasiswa)

	url := "https://script.google.com/macros/s/AKfycbztwrIURXKtvr2ldX8vNxhznQrxUWSLNPm8PuIveHHHkQu5R4ikv3pmMgel-NzptzE1wQ/exec"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		c.JSON(500, gin.H{"status": false, "pesan": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"status": false, "pesan": err.Error()})
		return
	}
	defer res.Body.Close()

	// Kalau Google kasih redirect, POST ulang ke URL redirect
	if res.StatusCode == http.StatusFound || res.StatusCode == http.StatusMovedPermanently {
		redirectURL := res.Header.Get("Location")

		req2, _ := http.NewRequest("POST", redirectURL, bytes.NewBuffer(postBody))
		req2.Header.Set("Content-Type", "application/json")

		res, err = http.DefaultClient.Do(req2)
		if err != nil {
			c.JSON(500, gin.H{"status": false, "pesan": err.Error()})
			return
		}
		defer res.Body.Close()
	}

	body, _ := io.ReadAll(res.Body)

	var hasil interface{}
	json.Unmarshal(body, &hasil)

	c.JSON(http.StatusOK, gin.H{
		"status_code": res.StatusCode,
		"kirim_data":  mahasiswa,
		"hasil_json":  hasil,
	})
}

func SheetUpdate(c *gin.Context) {
	var mahasiswa MahasiswaSheet

	if err := c.ShouldBindJSON(&mahasiswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})
		return
	}

	mahasiswa.Action = "update"

	postBody, _ := json.Marshal(mahasiswa)

	url := "https://script.google.com/macros/s/AKfycbztwrIURXKtvr2ldX8vNxhznQrxUWSLNPm8PuIveHHHkQu5R4ikv3pmMgel-NzptzE1wQ/exec"

	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(postBody),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  err.Error(),
		})
		return
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var hasil interface{}
	json.Unmarshal(body, &hasil)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"kirim_data": mahasiswa,
		"hasil_json": hasil,
	})
}
