package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v83"
)

type ClientInterface interface {
	CreateCustomer(ctx context.Context, params *stripe.CustomerCreateParams) (*stripe.Customer, error)
	GetCustomer(ctx context.Context, id string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error)
	CreateSubscription(ctx context.Context, params *stripe.SubscriptionCreateParams) (*stripe.Subscription, error)
}

type ClientWrapper struct {
	client *stripe.Client
}

func (w *ClientWrapper) CreateCustomer(ctx context.Context, params *stripe.CustomerCreateParams) (*stripe.Customer, error) {
	return w.client.V1Customers.Create(ctx, params)
}

func (w *ClientWrapper) GetCustomer(ctx context.Context, id string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error) {
	return w.client.V1Customers.Retrieve(ctx, id, params)
}

func (w *ClientWrapper) CreateSubscription(ctx context.Context, params *stripe.SubscriptionCreateParams) (*stripe.Subscription, error) {
	return w.client.V1Subscriptions.Create(ctx, params)
}
