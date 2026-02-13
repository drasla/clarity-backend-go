package maindb

import "gorm.io/gorm"

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type UserType string

const (
	TypePersonal UserType = "PERSONAL"
	TypeBusiness UserType = "BUSINESS"
)

type UserStatus string

const (
	StatusActive    UserStatus = "ACTIVE"
	StatusSuspended UserStatus = "SUSPENDED"
	StatusWithdrawn UserStatus = "WITHDRAWN"
)

type User struct {
	gorm.Model

	Username       string  `gorm:"type:varchar(50);uniqueIndex;not null;comment:사용자ID"`
	Password       string  `gorm:"type:varchar(255);not null;comment:암호화된 비밀번호"`
	Name           string  `gorm:"type:varchar(50);not null;comment:사용자 실명"`
	Email          string  `gorm:"type:varchar(100);uniqueIndex;not null;comment:이메일"`
	PhoneNumber    string  `gorm:"type:varchar(20);not null;comment:휴대폰번호"`
	LandlineNumber *string `gorm:"type:varchar(20);comment:일반전화번호"`

	Role   UserRole   `gorm:"type:varchar(20);default:'USER';not null;comment:권한(USER|ADMIN)"`
	Type   UserType   `gorm:"type:varchar(20);default:'PERSONAL';not null;comment:유형(PERSONAL|BUSINESS)"`
	Status UserStatus `gorm:"type:varchar(20);default:'PENDING';not null;index;comment:상태(ACTIVE|SUSPENDED|WITHDRAWN)"`

	AgreeEmail bool `gorm:"default:false;comment:이메일 수신 동의"`
	AgreeSMS   bool `gorm:"default:false;comment:SMS 수신 동의"`

	BizRegNumber  *string `gorm:"type:varchar(20);index;comment:사업자등록번호"`
	BizCEO        *string `gorm:"type:varchar(50);comment:대표자명"`
	BizType       *string `gorm:"type:varchar(50);comment:업태"`
	BizItem       *string `gorm:"type:varchar(50);comment:종목(업종)"`
	BizZipCode    *string `gorm:"type:varchar(10);comment:우편번호"`
	BizAddress1   *string `gorm:"type:varchar(255);comment:기본주소"`
	BizAddress2   *string `gorm:"type:varchar(255);comment:상세주소"`
	BizLicenseURL *string `gorm:"type:varchar(255);comment:사업자등록증 S3 URL"`
}

func (User) TableName() string {
	return "users"
}
