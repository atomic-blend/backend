package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v83"
)

type ClientInterface interface {
	CreateCustomer(ctx context.Context, params *stripe.CustomerCreateParams) (*stripe.Customer, error)
	GetCustomer(ctx context.Context, id string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error)
	CreateSubscription(ctx context.Context, params *stripe.SubscriptionCreateParams) (*stripe.Subscription, error)
	CreateInvoice(ctx context.Context, params *stripe.InvoiceCreateParams) (*stripe.Invoice, error)
	CreateInvoiceItem(ctx context.Context, params *stripe.InvoiceItemCreateParams) (*stripe.InvoiceItem, error)
	FinalizeInvoice(ctx context.Context, id string, params *stripe.InvoiceFinalizeInvoiceParams) (*stripe.Invoice, error)
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

func (w *ClientWrapper) CreateInvoice(ctx context.Context, params *stripe.InvoiceCreateParams) (*stripe.Invoice, error) {
	return w.client.V1Invoices.Create(ctx, params)
}

func (w *ClientWrapper) CreateInvoiceItem(ctx context.Context, params *stripe.InvoiceItemCreateParams) (*stripe.InvoiceItem, error) {
	return w.client.V1InvoiceItems.Create(ctx, params)
}

func (w *ClientWrapper) FinalizeInvoice(ctx context.Context, id string, params *stripe.InvoiceFinalizeInvoiceParams) (*stripe.Invoice, error) {
	return w.client.V1Invoices.FinalizeInvoice(ctx, id, params)
}