package maindb

import "gorm.io/gorm"

type EmailTemplate struct {
	gorm.Model
	TemplateCode string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Subject      string `gorm:"type:varchar(255);not null"`
	HTML         string `gorm:"type:longtext;not null;comment:발송용 컴파일된 HTML"`
	Design       string `gorm:"type:json;not null;comment:react-email-editor 디자인 JSON"`
	Variables    string `gorm:"type:text"`
	Description  string `gorm:"type:varchar(255)"`
}

func (EmailTemplate) TableName() string {
	return "email_templates"
}
