package models

import (
	"time"
)

// InviteCode represents an invitation code that grants access to Plex services
type InviteCode struct {
	ID              int        `json:"id"`
	Code            string     `json:"code"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	MaxUses         *int       `json:"max_uses,omitempty"`
	UsedCount       int        `json:"used_count"`
	IsDisabled      bool       `json:"is_disabled"`
	EntitlementName string     `json:"entitlement_name"`
	Duration        *time.Time `json:"duration,omitempty"`
}

// IsValid checks if an invite code is still valid for use
func (i *InviteCode) IsValid() bool {
	// Code is disabled
	if i.IsDisabled {
		return false
	}

	// Code has expired
	if i.ExpiresAt != nil && time.Now().After(*i.ExpiresAt) {
		return false
	}

	// Max uses reached
	if i.MaxUses != nil && i.UsedCount >= *i.MaxUses {
		return false
	}

	return true
}
