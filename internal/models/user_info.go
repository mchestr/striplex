package models

type UserInfo struct {
	ID       int     `json:"id"`
	UUID     string  `json:"uuid"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	IsAdmin  bool    `json:"is_admin"`
	Notes    *string `json:"notes,omitempty"` // Admin notes about the user, omitted if empty
}

type PlexAuth struct {
	PinID int    `json:"pin"`
	State string `json:"state"`
}
