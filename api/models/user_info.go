package models

type UserInfo struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PlexAuth struct {
	PinID int    `json:"pin"`
	State string `json:"state"`
}
