package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type validationError struct {
	FailedField string
	Tag         string
	Param       string
}

type XValidator struct {
	validator *validator.Validate
}

var validate = validator.New()
var BaseValidator = &XValidator{
	validator: validate,
}

func (v XValidator) Validate(data any) []validationError {
	var errs []validationError

	if err := v.validator.Struct(data); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				errs = append(errs, validationError{
					FailedField: e.Field(),
					Tag:         e.Tag(),
					Param:       e.Param(),
				})
			}
		}
	}

	return errs
}

func ValidateAPIData(data any) (bool, *fiber.Error) {
	if errs := BaseValidator.Validate(data); len(errs) > 0 {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msg := fmt.Sprintf("%s: failed '%s' validation", e.FailedField, e.Tag)
			if e.Param != "" {
				msg += fmt.Sprintf(" (must be %s)", e.Param)
			}
			msgs = append(msgs, msg)
		}
		return false, &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(msgs, "; "),
		}
	}
	return true, nil
}
