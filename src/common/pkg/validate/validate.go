package validate

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func Struct(s any) error {
	return validate.Struct(s)
}
