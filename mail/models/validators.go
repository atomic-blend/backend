package models

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterValidators registers custom validators for the models
func RegisterValidators() {
	if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
	}
}
