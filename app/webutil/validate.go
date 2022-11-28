package webutil

import "github.com/go-playground/validator/v10"

func InitValidator() *validator.Validate {
	Validate = validator.New()
	return Validate
}

var Validate *validator.Validate
