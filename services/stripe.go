package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"plefi/config"
	"plefi/model"
	"strconv"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
)

type StripeService struct {
	client *http.Client
}

func NewStripeService(client *http.Client) *StripeService {
	return &StripeService{
		client: client,
	}
}

func (s *StripeService) GetCustomer(ctx context.Context, user *model.UserInfo) (*stripe.Customer, error) {
	slog.Info("Searching for Stripe customer",
		"plex_id", user.ID,
		"email", user.Email,
		"username", user.Username)
	customerIter := customer.Search(&stripe.CustomerSearchParams{
		SearchParams: stripe.SearchParams{
			Query:   fmt.Sprintf("metadata['plex_user_id']:'%d'", user.ID),
			Limit:   stripe.Int64(1),
			Context: ctx,
		},
	})
	if err := customerIter.Err(); err != nil {
		return nil, err
	}
	if !customerIter.Next() {
		// No customer found create one
		return nil, nil
	}
	return customerIter.Customer(), nil
}

func (s *StripeService) GetOrCreateCustomer(ctx context.Context, user *model.UserInfo) (*stripe.Customer, error) {
	customer, err := s.GetCustomer(ctx, user)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return s.CreateCustomer(ctx, user)
	}
	return customer, nil
}

func (s *StripeService) CreateCustomer(ctx context.Context, user *model.UserInfo) (*stripe.Customer, error) {
	slog.Info("Creating a new Stripe customer",
		"plex_id", user.ID,
		"email", user.Email,
		"username", user.Username)

	// Create a Stripe customer first with Plex user metadata
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.Username),
		Metadata: map[string]string{
			"plex_user_id":  strconv.Itoa(user.ID),
			"plex_username": user.Username,
			"plex_email":    user.Email,
		},
		Params: stripe.Params{
			Context: ctx,
		},
	}
	return customer.New(customerParams)
}

func (s *StripeService) CreateSubscriptionCheckoutSession(ctx context.Context, sCustomer *stripe.Customer, user *model.UserInfo) (*stripe.CheckoutSession, error) {
	slog.Info("Creating a new Stripe subscription checkout session",
		"plex_id", user.ID,
		"email", user.Email,
		"username", user.Username)

	priceID := config.Config.GetString("stripe.default_price_id")
	hostname := config.Config.GetString("server.hostname")
	successURL := fmt.Sprintf("https://%s/stripe/success", hostname)
	cancelURL := fmt.Sprintf("https://%s/stripe/cancel?price_id=%s", hostname, priceID)
	// Create a Stripe checkout session for the customer
	return session.New(&stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Customer:   stripe.String(sCustomer.ID),
		Params: stripe.Params{
			Context: ctx,
		},
	})
}

func (s *StripeService) GetSubscription(ctx context.Context, userInfo *model.UserInfo, subscriptionID string) (*stripe.Subscription, error) {
	customer, err := s.GetCustomer(ctx, userInfo)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Validate subscription exists and belongs to this user
	subscription, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, err
	}
	// Verify subscription belongs to this user
	if subscription.Customer.ID != customer.ID {
		slog.Error("Subscription does not belong to this user",
			"plex_id", userInfo.ID,
			"email", userInfo.Email,
			"username", userInfo.Username,
			"subscription_id", subscription.ID,
			"customer_id", customer.ID)
		return nil, fmt.Errorf("subscription not found")
	}
	return subscription, nil
}

func (s *StripeService) CancelAtEndSubscription(ctx context.Context, subscriptionID string) (*stripe.Subscription, error) {
	// Cancel the subscription
	return subscription.Update(subscriptionID, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
		Params: stripe.Params{
			Context: ctx,
		},
	})
}
