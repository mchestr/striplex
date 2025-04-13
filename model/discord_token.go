package model

import "time"

type DiscordToken struct {
	AccessToken   string `gorm:"primaryKey"`
	RefreshToken  string
	ExpiresAt     time.Time
	Scopes        string
	DiscordUserId string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Status        int
}

func (d *DiscordToken) TableName() string {
	return "discord_token"
}
