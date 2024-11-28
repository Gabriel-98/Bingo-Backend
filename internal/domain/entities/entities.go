package entities

import (
	"time"
)

type User struct {
	Id        int64      `gorm:"type:bigserial;primaryKey"`
	Username  string     `gorm:"type:text;unique;not null"`
	Password  string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;autoUpdateTime"`
}

type RefreshToken struct {
	Token string         `gorm:"type:string;primaryKey"`
	UserId int64         `gorm:"type:bigint;not null"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null"`
	ExpiresAt time.Time  `gorm:"type:timestamptz;not null"`
}