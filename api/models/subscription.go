package models

import "github.com/stripe/stripe-go/v82"

// SimplifiedSubscription represents minimal subscription data needed by frontend
type SubscriptionSummary struct {
	CustomerID        string             `json:"customer_id"`
	ID                string             `json:"id"`
	Status            string             `json:"status"`
	CancelAtPeriodEnd bool               `json:"cancel_at_period_end"`
	CancelAt          int64              `json:"cancel_at"`
	Items             []SubscriptionItem `json:"items"`
}

// SubscriptionItem holds minimal item data
type SubscriptionItem struct {
	ID               string    `json:"id"`
	PriceID          string    `json:"price_id"`
	Quantity         int64     `json:"quantity"`
	CurrentPeriodEnd int64     `json:"current_period_end"`
	PriceItem        PriceItem `json:"price"`
}

type PriceItem struct {
	UnitAmount int64         `json:"unit_amount"`
	Currency   string        `json:"currency"`
	Recurring  RecurringItem `json:"recurring"`
}

type RecurringItem struct {
	Interval string `json:"interval"`
}

// NewSubscriptionSummary maps a Stripe subscription to our minimal model.
func NewSubscriptionSummary(s *stripe.Subscription) *SubscriptionSummary {
	if s == nil {
		return nil
	}
	items := make([]SubscriptionItem, 0, len(s.Items.Data))
	for _, it := range s.Items.Data {
		items = append(items, SubscriptionItem{
			ID:               it.ID,
			PriceID:          it.Price.ID,
			Quantity:         it.Quantity,
			CurrentPeriodEnd: it.CurrentPeriodEnd,
			PriceItem: PriceItem{
				UnitAmount: it.Price.UnitAmount,
				Currency:   string(it.Price.Currency),
				Recurring: RecurringItem{
					Interval: string(it.Price.Recurring.Interval),
				},
			},
		})
	}
	return &SubscriptionSummary{
		CustomerID:        s.Customer.ID,
		ID:                s.ID,
		Status:            string(s.Status),
		CancelAtPeriodEnd: s.CancelAtPeriodEnd,
		CancelAt:          s.CancelAt,
		Items:             items,
	}
}
