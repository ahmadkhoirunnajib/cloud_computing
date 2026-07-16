package models

type Pesan struct {
	ID      uint   `gorm:"primaryKey"`
	Kode    string `gorm:"unique"`
	Balasan string `gorm:"type:text"`
}

func (Pesan) TableName() string {
	return "pesan"
}
