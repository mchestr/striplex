package model

import "time"

type PlexUser struct {
	ID            string `gorm:"primaryKey"`
	Username      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DiscordUserId string
	IsSubscriber  bool
	PlexTokens   []PlexToken `gorm:"foreignKey:PlexUserId"`
}

func (d *PlexUser) TableName() string {
	return "plex_user"
}
