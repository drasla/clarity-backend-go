package maindb

import (
	"time"

	"gorm.io/gorm"
)

type VerificationType string

const (
	VerifyEmail VerificationType = "EMAIL"
	VerifySMS   VerificationType = "SMS"
)

type Verification struct {
	gorm.Model

	Target     string           `gorm:"type:varchar(100);index;not null;comment:이메일 또는 전화번호"`
	Type       VerificationType `gorm:"type:varchar(10);not null;comment:인증유형(EMAIL|SMS)"`
	Code       string           `gorm:"type:varchar(10);not null;comment:인증코드(6자리)"`
	IsVerified bool             `gorm:"default:false;comment:인증완료여부"`
	ExpiresAt  time.Time        `gorm:"not null;index;comment:만료시간"`

	AttemptCount int `gorm:"default:1;comment:시도횟수"`
}

func (Verification) TableName() string {
	return "verifications"
}
