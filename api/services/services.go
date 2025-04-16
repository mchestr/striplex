package services

import "net/http"

type Services struct {
	Plex   PlexServicer
	Stripe StripeServicer
}

func NewServices(client *http.Client) *Services {
	return &Services{
		Plex:   NewPlexService(client),
		Stripe: NewStripeService(client),
	}
}
