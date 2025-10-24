package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserEntity represents a user in the system
// @Summary User entity
// @Description Represents a user in the system
type UserEntity struct {
	ID                       *primitive.ObjectID   `json:"id" bson:"_id"`
	StripeCustomerID         *string               `json:"stripeCustomerId" bson:"stripe_customer_id"`
	StripeSubscriptionID     *string               `json:"stripeSubscriptionId" bson:"stripe_subscription_id"`
	StripeSubscriptionStatus *string               `json:"stripeSubscriptionStatus" bson:"stripe_subscription_status"`
	Email                    *string               `json:"email" bson:"email" binding:"required"`
	FirstName                *string               `json:"firstName" bson:"first_name,omitempty"`
	LastName                 *string               `json:"lastName" bson:"last_name,omitempty"`
	BackupEmail              *string               `json:"backupEmail" bson:"backup_email"`
	Password                 *string               `json:"password,omitempty" bson:"password" binding:"required"`
	KeySet                   *EncryptionKey        `json:"keySet,omitempty" bson:"key_set" binding:"required"`
	RoleIds                  []*primitive.ObjectID `json:"-" bson:"role_ids"`
	Roles                    []*UserRoleEntity     `json:"roles" bson:"roles,omitempty"`
	ResetPasswordCode        *string               `json:"resetPasswordCode,omitempty" bson:"reset_password_code"`
	Devices                  []*UserDevice         `json:"devices" bson:"devices,omitempty"`
	Purchases                []*PurchaseEntity     `json:"purchases" bson:"purchases,omitempty"`
	CreatedAt                *primitive.DateTime   `json:"createdAt" bson:"created_at"`
	UpdatedAt                *primitive.DateTime   `json:"updatedAt" bson:"updated_at"`
}
