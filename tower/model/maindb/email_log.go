package maindb

import "gorm.io/gorm"

type EmailLogStatus string

const (
	EmailLogStatusSuccess EmailLogStatus = "SUCCESS"
	EmailLogStatusFailed  EmailLogStatus = "FAILED"
)

type EmailLog struct {
	gorm.Model

	TemplateCode string         `gorm:"type:varchar(50);index"`
	Sender       string         `gorm:"type:varchar(255);not null"`
	Recipient    string         `gorm:"type:varchar(255);not null;index"`
	Subject      string         `gorm:"type:varchar(255);not null"`
	HTMLBody     string         `gorm:"type:text;not null"`
	Status       EmailLogStatus `gorm:"type:varchar(20);not null;index"`
	ProviderID   string         `gorm:"type:varchar(255)"`
	ErrorMessage string         `gorm:"type:text"`
}
