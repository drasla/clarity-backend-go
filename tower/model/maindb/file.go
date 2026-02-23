package maindb

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	OriginalName string `gorm:"type:varchar(255);not null"`
	StoredName   string `gorm:"type:varchar(255);not null"`
	URL          string `gorm:"type:varchar(500);not null"`
	Size         int64  `gorm:"not null"`
	Extension    string `gorm:"type:varchar(50)"`
	TargetID     uint   `gorm:"index"`
	TargetType   string `gorm:"type:varchar(50);index"`
}
