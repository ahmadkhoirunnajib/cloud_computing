package models

import "time"

type Suhu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Lokasi    string    `json:"lokasi"`
	Suhu      float32   `json:"suhu"`
	CreatedAt time.Time `json:"created_at"`
}
