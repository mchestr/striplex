package services

import (
	"net/http"
	"plefi/internal/services/plex"
)

type Services struct {
	Plex   plex.PlexServicer
	Stripe StripeServicer
}

func NewServices(client *http.Client) *Services {
	return &Services{
		Plex:   plex.NewPlexService(client),
		Stripe: NewStripeService(client),
	}
}
