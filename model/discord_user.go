package model

import "time"

type DiscordUser struct {
	ID            string `gorm:"primaryKey"`
	Username      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsActive      bool
	DiscordTokens []DiscordToken `gorm:"foreignKey:DiscordUserId"`
	PlexUsers     []PlexUser     `gorm:"foreignKey:DiscordUserId"`
}

func (d *DiscordUser) TableName() string {
	return "discord_user"
}
