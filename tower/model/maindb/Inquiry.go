package maindb

import (
	"time"

	"gorm.io/gorm"
)

type InquiryStatus string

const (
	InquiryStatusPending   InquiryStatus = "PENDING"
	InquiryStatusCompleted InquiryStatus = "COMPLETED"
)

type InquiryCategory string

const (
	InquiryCategoryDomain     InquiryCategory = "DOMAIN"
	InquiryCategoryHosting    InquiryCategory = "HOSTING"
	InquiryCategoryGoldenShop InquiryCategory = "GOLDEN_SHOP"
	InquiryCategoryEmail      InquiryCategory = "EMAIL"
	InquiryCategorySSL        InquiryCategory = "SSL"
	InquiryCategoryUserInfo   InquiryCategory = "USER_INFO"
	InquiryCategoryEtc        InquiryCategory = "ETC"
)

type Inquiry struct {
	gorm.Model

	UserID      *uint           `gorm:"index"`
	NonMemberPw *string         `gorm:"type:varchar(255);comment:비회원 비밀번호"`
	Category    InquiryCategory `gorm:"type:varchar(50);not null"`
	Domain      *string         `gorm:"type:varchar(255)"`
	Title       string          `gorm:"type:varchar(255);not null"`
	Content     string          `gorm:"type:text;not null"`
	Email       string          `gorm:"type:varchar(255);not null"`
	PhoneNumber string          `gorm:"type:varchar(50);not null"`
	Status      InquiryStatus   `gorm:"type:varchar(20);default:'PENDING'"`
	Answer      *string         `gorm:"type:text"`
	AnsweredAt  *time.Time

	Attachments []File `gorm:"polymorphic:Target;"`
}
