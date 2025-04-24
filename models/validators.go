package models

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterValidators registers custom validators for the models
func RegisterValidators() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("validFrequency", ValidateFrequency)
        v.RegisterValidation("validPurchaseType", ValidatePurchaseType)
    }
}

// ValidateFrequency checks if the frequency value is in the allowed list
func ValidateFrequency(fl validator.FieldLevel) bool {
    frequency := fl.Field().String()
    for _, validFreq := range ValidFrequencies {
        if frequency == validFreq {
            return true
        }
    }
    return false
}

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