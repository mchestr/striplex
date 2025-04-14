package services

import "net/http"

type Services struct {
	Wizarr *WizarrService
	Plex   *PlexService
}

func NewServices(client *http.Client) *Services {
	return &Services{
		Wizarr: NewWizarrService(client),
		Plex:   NewPlexService(client),
	}
}
