// Package stripe provides a wrapper around the Stripe client to facilitate testing and abstraction.
package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v83"
)

// ClientInterface defines the methods our Stripe client wrapper must implement.
type ClientInterface interface {
	CreateCustomer(ctx context.Context, params *stripe.CustomerCreateParams) (*stripe.Customer, error)
	GetCustomer(ctx context.Context, id string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error)
	CreateSubscription(ctx context.Context, params *stripe.SubscriptionCreateParams) (*stripe.Subscription, error)
	CreateInvoice(ctx context.Context, params *stripe.InvoiceCreateParams) (*stripe.Invoice, error)
	CreateInvoiceItem(ctx context.Context, params *stripe.InvoiceItemCreateParams) (*stripe.InvoiceItem, error)
	FinalizeInvoice(ctx context.Context, id string, params *stripe.InvoiceFinalizeInvoiceParams) (*stripe.Invoice, error)
	GetEphemeralKeys(ctx context.Context, params *stripe.EphemeralKeyCreateParams) (*stripe.EphemeralKey, error)
	CreateCheckoutSession(ctx context.Context, params *stripe.CheckoutSessionCreateParams) (*stripe.CheckoutSession, error)
}

// ClientWrapper wraps the Stripe client to implement ClientInterface.
type ClientWrapper struct {
	client *stripe.Client
}

// CreateCustomer creates a new Stripe customer.
func (w *ClientWrapper) CreateCustomer(ctx context.Context, params *stripe.CustomerCreateParams) (*stripe.Customer, error) {
	return w.client.V1Customers.Create(ctx, params)
}

// GetCustomer retrieves a Stripe customer by ID.
func (w *ClientWrapper) GetCustomer(ctx context.Context, id string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error) {
	return w.client.V1Customers.Retrieve(ctx, id, params)
}

// CreateSubscription creates a new Stripe subscription.
func (w *ClientWrapper) CreateSubscription(ctx context.Context, params *stripe.SubscriptionCreateParams) (*stripe.Subscription, error) {
	return w.client.V1Subscriptions.Create(ctx, params)
}

// CreateInvoice creates a new Stripe invoice.
func (w *ClientWrapper) CreateInvoice(ctx context.Context, params *stripe.InvoiceCreateParams) (*stripe.Invoice, error) {
	return w.client.V1Invoices.Create(ctx, params)
}

// CreateInvoiceItem creates a new Stripe invoice item.
func (w *ClientWrapper) CreateInvoiceItem(ctx context.Context, params *stripe.InvoiceItemCreateParams) (*stripe.InvoiceItem, error) {
	return w.client.V1InvoiceItems.Create(ctx, params)
}

// FinalizeInvoice finalizes a Stripe invoice.
func (w *ClientWrapper) FinalizeInvoice(ctx context.Context, id string, params *stripe.InvoiceFinalizeInvoiceParams) (*stripe.Invoice, error) {
	return w.client.V1Invoices.FinalizeInvoice(ctx, id, params)
}

// GetEphemeralKeys creates a new Stripe ephemeral key.
func (w *ClientWrapper) GetEphemeralKeys(ctx context.Context, params *stripe.EphemeralKeyCreateParams) (*stripe.EphemeralKey, error) {
	return w.client.V1EphemeralKeys.Create(ctx, params)
}

func (w *ClientWrapper) CreateCheckoutSession(ctx context.Context, params *stripe.CheckoutSessionCreateParams) (*stripe.CheckoutSession, error) {
	return w.client.V1CheckoutSessions.Create(ctx, params)
}
