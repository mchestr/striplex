package model

import "time"

type PlexToken struct {
	AccessToken string `gorm:"primaryKey"`
	PlexUserId  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *PlexToken) TableName() string {
	return "plex_token"
}
