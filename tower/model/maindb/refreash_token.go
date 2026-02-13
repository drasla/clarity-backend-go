package maindb

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model

	UserID    uint      `gorm:"not null;index;comment:User 테이블 PK"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null;comment:리프레시 토큰 값"`
	ExpiresAt time.Time `gorm:"not null;comment:토큰 만료 시간"`
	IsRevoked bool      `gorm:"default:false;comment:토큰 폐기 여부(강제 로그아웃용)"`
	ClientIP  string    `gorm:"type:varchar(50);comment:로그인한 IP 주소"`
	UserAgent string    `gorm:"type:varchar(255);comment:로그인한 기기/브라우저 정보"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
