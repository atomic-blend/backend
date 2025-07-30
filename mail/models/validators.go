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


// ValidatePurchaseType checks if the purchase type is in the allowed list
func ValidatePurchaseType(fl validator.FieldLevel) bool {
    purchaseType := fl.Field().String()
    for _, validType := range ValidPurchaseTypes {
        if purchaseType == validType {
            return true
        }
    }
    return false
}