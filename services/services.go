package services

import "net/http"

type Services struct {
	Plex   *PlexService
	Stripe *StripeService
}

func NewServices(client *http.Client) *Services {
	return &Services{
		Plex:   NewPlexService(client),
		Stripe: NewStripeService(client),
	}
}
