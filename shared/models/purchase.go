package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ValidPurchaseTypes is used in validators.go
var ValidPurchaseTypes = []string{
	"REVENUE_CAT",
	"STRIPE",
}

// PurchaseEntity represents a generic purchase record using a tagged union pattern
// Only one of RevenueCatData or StripeData should be non-nil, indicated by the Type field
type PurchaseEntity struct {
	ID             primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	Type           *string                 `json:"type" bson:"type" binding:"required,validPurchaseType"`
	RevenueCatData *RevenueCatPurchaseData `json:"revenueCatData,omitempty" bson:"revenue_cat_data,omitempty"`
	StripeData     *StripePurchaseData     `json:"stripeData,omitempty" bson:"stripe_data,omitempty"`
	CreatedAt      primitive.DateTime      `json:"createdAt" bson:"created_at"`
	UpdatedAt      primitive.DateTime      `json:"updatedAt" bson:"updated_at"`
}

// RevenueCatPurchaseData represents RevenueCat specific purchase data
type RevenueCatPurchaseData struct {
	Aliases                  []string             `json:"aliases"`
	AppID                    string               `json:"app_id"`
	AppUserID                string               `json:"app_user_id"`
	CommissionPercentage     float64              `json:"commission_percentage"`
	CountryCode              string               `json:"country_code"`
	Currency                 string               `json:"currency"`
	EntitlementID            string               `json:"entitlement_id"`
	EntitlementIDs           []string             `json:"entitlement_ids"`
	Environment              string               `json:"environment"`
	EventTimestampMs         int64                `json:"event_timestamp_ms"`
	ExpirationAtMs           int64                `json:"expiration_at_ms"`
	ID                       string               `json:"id"`
	IsFamilyShare            bool                 `json:"is_family_share"`
	OfferCode                string               `json:"offer_code"`
	OriginalAppUserID        string               `json:"original_app_user_id"`
	OriginalTransactionID    string               `json:"original_transaction_id"`
	PeriodType               string               `json:"period_type"`
	PresentedOfferingID      string               `json:"presented_offering_id"`
	Price                    float64              `json:"price"`
	PriceInPurchasedCurrency float64              `json:"price_in_purchased_currency"`
	ProductID                string               `json:"product_id"`
	PurchasedAtMs            int64                `json:"purchased_at_ms"`
	Store                    string               `json:"store"`
	SubscriberAttributes     map[string]Attribute `json:"subscriber_attributes"`
	TakehomePercentage       float64              `json:"takehome_percentage"`
	TaxPercentage            float64              `json:"tax_percentage"`
	TransactionID            string               `json:"transaction_id"`
	Type                     string               `json:"type"`
}

// Attribute represents a subscriber attribute in RevenueCat
type Attribute struct {
	UpdatedAtMs int64  `json:"updated_at_ms"`
	Value       string `json:"value"`
}

// RevenueCatPayload represents the structure of a RevenueCat webhook payload
type RevenueCatPayload struct {
	APIVersion string                 `json:"api_version"`
	Event      RevenueCatPurchaseData `json:"event"`
}

// StripePurchaseData represents Stripe specific purchase data
type StripePurchaseData struct {
	ID             string            `json:"id"`
	Object         string            `json:"object"`
	Amount         int64             `json:"amount"`
	AmountCaptured int64             `json:"amount_captured"`
	AmountRefunded int64             `json:"amount_refunded"`
	Currency       string            `json:"currency"`
	CustomerID     string            `json:"customer"`
	Description    string            `json:"description"`
	Invoice        string            `json:"invoice"`
	Metadata       map[string]string `json:"metadata"`
	PaymentIntent  string            `json:"payment_intent"`
	PaymentMethod  string            `json:"payment_method"`
	ReceiptEmail   string            `json:"receipt_email"`
	ReceiptURL     string            `json:"receipt_url"`
	Status         string            `json:"status"`
	SubscriptionID string            `json:"subscription"`
	Created        int64             `json:"created"`
	// Add additional Stripe-specific fields as needed
}

// StripePayload represents the structure of a Stripe webhook payload
type StripePayload struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Type    string          `json:"type"`
	Data    StripeEventData `json:"data"`
	Created int64           `json:"created"`
}

// StripeEventData represents the data field in a Stripe webhook event
type StripeEventData struct {
	Object StripePurchaseData `json:"object"`
}

// NewRevenueCatPurchase creates a new PurchaseEntity with RevenueCat data
func NewRevenueCatPurchase(rcData RevenueCatPurchaseData) PurchaseEntity {
	now := primitive.NewDateTimeFromTime(time.Now())
	return PurchaseEntity{
		ID:             primitive.NewObjectID(),
		Type:           stringPtr("REVENUE_CAT"),
		RevenueCatData: &rcData,
		StripeData:     nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// NewStripePurchase creates a new PurchaseEntity with Stripe data
func NewStripePurchase(stripeData StripePurchaseData) PurchaseEntity {
	now := primitive.NewDateTimeFromTime(time.Now())
	return PurchaseEntity{
		ID:             primitive.NewObjectID(),
		Type:           stringPtr("STRIPE"),
		RevenueCatData: nil,
		StripeData:     &stripeData,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// IsRevenueCat returns true if this purchase is from RevenueCat
func (p *PurchaseEntity) IsRevenueCat() bool {
	return p.Type != nil && *p.Type == "REVENUE_CAT" && p.RevenueCatData != nil
}

// IsStripe returns true if this purchase is from Stripe
func (p *PurchaseEntity) IsStripe() bool {
	return p.Type != nil && *p.Type == "STRIPE" && p.StripeData != nil
}

// GetRevenueCatData safely retrieves the RevenueCat data if available
func (p *PurchaseEntity) GetRevenueCatData() (*RevenueCatPurchaseData, bool) {
	if p.IsRevenueCat() {
		return p.RevenueCatData, true
	}
	return nil, false
}

// GetStripeData safely retrieves the Stripe data if available
func (p *PurchaseEntity) GetStripeData() (*StripePurchaseData, bool) {
	if p.IsStripe() {
		return p.StripeData, true
	}
	return nil, false
}

// stringPtr is a helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}
