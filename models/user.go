package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `json:"nama"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
