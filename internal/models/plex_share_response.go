package models

import "time"

// PlexShareResponse represents the JSON returned by the Plex share API
type PlexShareResponse struct {
	Accepted     bool       `json:"accepted"`
	AcceptedAt   *time.Time `json:"acceptedAt"`
	AllLibraries bool       `json:"allLibraries"`
	DeletedAt    *time.Time `json:"deletedAt"`
	ID           int        `json:"id"`
	InviteToken  string     `json:"inviteToken"`
	Invited      struct {
		FriendlyName    *string `json:"friendlyName"`
		Home            bool    `json:"home"`
		ID              int     `json:"id"`
		Restricted      bool    `json:"restricted"`
		SharingSettings struct {
			AllowCameraUpload  bool    `json:"allowCameraUpload"`
			AllowChannels      bool    `json:"allowChannels"`
			AllowSubtitleAdmin bool    `json:"allowSubtitleAdmin"`
			AllowSync          bool    `json:"allowSync"`
			AllowTuners        int     `json:"allowTuners"`
			FilterAll          *string `json:"filterAll"`
			FilterMovies       string  `json:"filterMovies"`
			FilterMusic        string  `json:"filterMusic"`
			FilterPhotos       string  `json:"filterPhotos"`
			FilterTelevision   string  `json:"filterTelevision"`
		} `json:"sharingSettings"`
		Status   string `json:"status"`
		Thumb    string `json:"thumb"`
		Title    string `json:"title"`
		Username string `json:"username"`
		UUID     string `json:"uuid"`
	} `json:"invited"`
	InvitedEmail *string    `json:"invitedEmail"`
	InvitedID    int        `json:"invitedId"`
	LastSeenAt   time.Time  `json:"lastSeenAt"`
	LeftAt       *time.Time `json:"leftAt"`
	Libraries    []struct {
		ID    int    `json:"id"`
		Key   int    `json:"key"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"libraries"`
	MachineIdentifier string `json:"machineIdentifier"`
	Name              string `json:"name"`
	NumLibraries      int    `json:"numLibraries"`
	Owned             bool   `json:"owned"`
	OwnerID           int    `json:"ownerId"`
	ServerID          int    `json:"serverId"`
	SharingSettings   struct {
		AllowCameraUpload  bool    `json:"allowCameraUpload"`
		AllowChannels      bool    `json:"allowChannels"`
		AllowSubtitleAdmin bool    `json:"allowSubtitleAdmin"`
		AllowSync          bool    `json:"allowSync"`
		AllowTuners        int     `json:"allowTuners"`
		FilterAll          *string `json:"filterAll"`
		FilterMovies       string  `json:"filterMovies"`
		FilterMusic        string  `json:"filterMusic"`
		FilterPhotos       string  `json:"filterPhotos"`
		FilterTelevision   string  `json:"filterTelevision"`
	} `json:"sharingSettings"`
}
