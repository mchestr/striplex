package model

import "time"

type DiscordUser struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}
