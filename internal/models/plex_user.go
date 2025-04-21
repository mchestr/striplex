package models

import (
	"time"
)

// PlexUser represents user information from Plex
type PlexUser struct {
	ID        int       `json:"id"`         // Plex user ID
	UUID      string    `json:"uuid"`       // Plex user UUID
	Username  string    `json:"username"`   // Plex username
	Email     string    `json:"email"`      // Plex email
	IsAdmin   bool      `json:"is_admin"`   // Is this user an admin
	CreatedAt time.Time `json:"created_at"` // When the user was created in our system
	UpdatedAt time.Time `json:"updated_at"` // When the user was last updated in our system
}

// PlexUserInvite associates a user with an invite code they've used
type PlexUserInvite struct {
	ID              int        `json:"id"`               // Primary key
	UserID          int        `json:"user_id"`          // Plex user ID
	InviteCodeID    int        `json:"invite_code_id"`   // Invite code ID they used
	InviteCode      string     `json:"invite_code"`      // The actual code (populated from join)
	EntitlementName string     `json:"entitlement_name"` // Entitlement from the invite code
	UsedAt          time.Time  `json:"used_at"`          // When the code was used
	ExpiresAt       *time.Time `json:"expires_at"`       // When the access expires (if applicable)
}

// HasValidAccess checks if the invite still provides valid access
func (pui *PlexUserInvite) HasValidAccess() bool {
	// If there's no expiration, it's valid forever
	if pui.ExpiresAt == nil {
		return true
	}
	// Otherwise check if it's expired
	return time.Now().Before(*pui.ExpiresAt)
}
