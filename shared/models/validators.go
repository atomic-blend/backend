package models

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterValidators registers custom validators for the models
func RegisterValidators() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("validPurchaseType", ValidatePurchaseType)
    }
}

// ValidatePurchaseType is a list of valid frequency values for subscriptions
func ValidatePurchaseType(fl validator.FieldLevel) bool {
    purchaseType := fl.Field().String()
    validPurchaseTypes := []string{}
    for _, validType := range validPurchaseTypes {
        if purchaseType == validType {
            return true
        }
    }
    return false
}