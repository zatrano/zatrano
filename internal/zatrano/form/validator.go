package form

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		if name == "-" { return "" }
		if name == "" { return fld.Name }
		return name
	})
}

func Validate(s interface{}) (map[string]string, error) {
	err := validate.Struct(s)
	if err == nil { return nil, nil }
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok { return nil, err }
	errorMap := make(map[string]string)
	for _, e := range validationErrors {
		errorMap[e.Field()] = formatErrorMessage(e)
	}
	return errorMap, err
}

func formatErrorMessage(e validator.FieldError) string {
	fieldName := e.Field()
	switch e.Tag() {
	case "required": return fmt.Sprintf("The %s field is required.", fieldName)
	case "email": return "Please provide a valid email address."
	case "min": return fmt.Sprintf("The %s field must be at least %s characters long.", fieldName, e.Param())
	case "max": return fmt.Sprintf("The %s field must be at most %s characters long.", fieldName, e.Param())
	case "eqfield": return fmt.Sprintf("The %s field must match the %s field.", fieldName, e.Param())
	default: return fmt.Sprintf("The %s field is not valid.", fieldName)
	}
}