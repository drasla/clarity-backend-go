package maindb

import "gorm.io/gorm"

type EmailTemplate struct {
	gorm.Model
	TemplateCode string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Subject      string `gorm:"type:varchar(255);not null"`
	HTMLBody     string `gorm:"type:text;not null"`
	Variables    string `gorm:"type:text"`
	Description  string `gorm:"type:varchar(255)"`
}
