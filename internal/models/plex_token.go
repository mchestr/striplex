package models

import "time"

// PlexToken holds a user's Plex OAuth tokens.
type PlexToken struct {
	UserID      int       `db:"user_id"    json:"user_id"`
	AccessToken string    `db:"access_token"  json:"access_token"`
	CreatedAt   time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"    json:"updated_at"`
}
