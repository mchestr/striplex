package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"plefi/api/config"
	"plefi/api/models"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
)

// StripeServicer defines the interface for Stripe payment operations
type StripeServicer interface {
	// GetCustomer retrieves an existing Stripe customer by Plex user ID
	GetCustomer(ctx context.Context, user *models.UserInfo) (*stripe.Customer, error)

	// GetOrCreateCustomer retrieves a customer or creates one if it doesn't exist
	GetOrCreateCustomer(ctx context.Context, user *models.UserInfo) (*stripe.Customer, error)

	// CreateCustomer creates a new Stripe customer from Plex user info
	CreateCustomer(ctx context.Context, user *models.UserInfo) (*stripe.Customer, error)

	// CreateAnonymousCustomer creates a customer for anonymous donations
	CreateAnonymousCustomer(ctx context.Context) (*stripe.Customer, error)

	// CreateSubscriptionCheckoutSession creates a checkout session for subscription purchase
	CreateSubscriptionCheckoutSession(ctx context.Context, sCustomer *stripe.Customer, user *models.UserInfo, anchorDate *time.Time) (*stripe.CheckoutSession, error)

	// CreateOneTimeCheckoutSession creates a checkout session for one-time payment
	CreateOneTimeCheckoutSession(ctx context.Context, sCustomer *stripe.Customer, user *models.UserInfo) (*stripe.CheckoutSession, error)

	// GetSubscription retrieves a subscription and verifies it belongs to the user
	GetSubscription(ctx context.Context, userInfo *models.UserInfo, subscriptionID string) (*models.SubscriptionSummary, error)

	// CancelAtEndSubscription cancels a subscription at the end of the current period
	CancelAtEndSubscription(ctx context.Context, subscriptionID string) (*models.SubscriptionSummary, error)

	// GetActiveSubscriptions returns all active subscriptions for a user
	GetActiveSubscription(ctx context.Context, user *models.UserInfo) (*models.SubscriptionSummary, error)
}

// Verify that StripeService implements the StripeServicer interface
var _ StripeServicer = (*StripeService)(nil)

type StripeService struct {
	client *http.Client
}

func NewStripeService(client *http.Client) *StripeService {
	return &StripeService{
		client: client,
	}
}

func (s *StripeService) GetCustomer(ctx context.Context, user *models.UserInfo) (*stripe.Customer, error) {
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

func (s *StripeService) GetOrCreateCustomer(ctx context.Context, user *models.UserInfo) (*stripe.Customer, error) {
	customer, err := s.GetCustomer(ctx, user)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return s.CreateCustomer(ctx, user)
	}
	return customer, nil
}

func (s *StripeService) CreateCustomer(ctx context.Context, user *models.UserInfo) (*stripe.Customer, error) {
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

// CreateAnonymousCustomer creates a customer for anonymous donations
func (s *StripeService) CreateAnonymousCustomer(ctx context.Context) (*stripe.Customer, error) {
	slog.Info("Creating an anonymous Stripe customer for donation")

	return customer.New(&stripe.CustomerParams{
		Description: stripe.String("Anonymous donation customer"),
		Params: stripe.Params{
			Context: ctx,
		},
	})
}

func (s *StripeService) CreateSubscriptionCheckoutSession(ctx context.Context, sCustomer *stripe.Customer, user *models.UserInfo, anchorDate *time.Time) (*stripe.CheckoutSession, error) {
	slog.Info("Creating a new Stripe subscription checkout session",
		"plex_id", user.ID,
		"email", user.Email,
		"username", user.Username)

	successURL := fmt.Sprintf("https://%s/stripe/success", config.C.Server.Hostname)
	cancelURL := fmt.Sprintf("https://%s/stripe/cancel", config.C.Server.Hostname)
	// Create a Stripe checkout session for the customer
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(config.C.Stripe.SubscriptionPriceID),
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
	}
	if anchorDate != nil {
		slog.Info("setting anchor date", "anchor_date", anchorDate.Format(time.RFC3339))
		params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{
			TrialEnd: stripe.Int64(anchorDate.Unix()),
		}
	}
	return session.New(params)
}

func (s *StripeService) CreateOneTimeCheckoutSession(ctx context.Context, sCustomer *stripe.Customer, user *models.UserInfo) (*stripe.CheckoutSession, error) {
	// Log with user info if available
	if user != nil {
		slog.Info("Creating a new Stripe donation checkout session",
			"plex_id", user.ID,
			"email", user.Email,
			"username", user.Username)
	} else {
		slog.Info("Creating an anonymous donation checkout session")
	}

	successURL := fmt.Sprintf("https://%s/stripe/donation-success", config.C.Server.Hostname)
	cancelURL := fmt.Sprintf("https://%s/stripe/cancel", config.C.Server.Hostname)
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice(config.C.Stripe.PaymentMethodTypes),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(config.C.Stripe.DonationPriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Params: stripe.Params{
			Context: ctx,
		},
	}
	if sCustomer != nil {
		params.Customer = stripe.String(sCustomer.ID)
	}

	// Create a Stripe checkout session for the customer
	return session.New(params)
}

func (s *StripeService) GetSubscription(ctx context.Context, userInfo *models.UserInfo, subscriptionID string) (*models.SubscriptionSummary, error) {
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
	return models.NewSubscriptionSummary(subscription), nil
}

func (s *StripeService) GetActiveSubscription(ctx context.Context, user *models.UserInfo) (*models.SubscriptionSummary, error) {
	customer, err := s.GetCustomer(ctx, user)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	iter := subscription.List(&stripe.SubscriptionListParams{
		Customer: stripe.String(customer.ID),
		ListParams: stripe.ListParams{
			Context: ctx,
		},
	})
	if err := iter.Err(); err != nil {
		return nil, err
	}

	var sub *stripe.Subscription
	for iter.Next() {
		sub = iter.Subscription()
		if !(sub.Status == stripe.SubscriptionStatusActive || sub.Status == stripe.SubscriptionStatusTrialing) {
			continue
		}
		if !sub.CancelAtPeriodEnd {
			return models.NewSubscriptionSummary(sub), nil
		}
	}
	return models.NewSubscriptionSummary(sub), nil
}

func (s *StripeService) CancelAtEndSubscription(ctx context.Context, subscriptionID string) (*models.SubscriptionSummary, error) {
	// Cancel the subscription
	sub, err := subscription.Update(subscriptionID, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
		Params: stripe.Params{
			Context: ctx,
		},
	})
	if err != nil {
		return nil, err
	}
	return models.NewSubscriptionSummary(sub), nil
}
