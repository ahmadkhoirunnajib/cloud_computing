package models

import "gorm.io/gorm"

type Dokumen struct {
	gorm.Model

	NamaDokumen string `json:"nama_dokumen"`
	FileId      string `json:"file_id"`
	FileUrl     string `json:"file_url"`
}
