package models

import "gorm.io/gorm"

type Nilai struct {
	gorm.Model

	NamaMhs    string  `json:"nama_mhs"`
	NilaiTugas float64 `json:"nilai_tugas"`
	NilaiUts   float64 `json:"nilai_uts"`
	NilaiUas   float64 `json:"nilai_uas"`

	NilaiAngka float64 `json:"nilai_angka"`
	NilaiHuruf string  `json:"nilai_huruf"`
	Status     string  `json:"status"`
}
