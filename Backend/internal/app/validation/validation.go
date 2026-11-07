package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func Validate(s any) error {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return errors.New(formatErrors(ve))
	}
	return err
}

func formatErrors(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "ошибка валидации"
	}

	var parts []string
	for _, e := range errs {
		parts = append(parts, fmt.Sprintf("%s: %s", fieldName(e), ruleMessage(e)))
	}
	return strings.Join(parts, "; ")
}

func fieldName(e validator.FieldError) string {
	field := e.StructField()
	if tag := e.StructNamespace(); tag != "" {
		if idx := strings.LastIndex(tag, "."); idx != -1 {
			field = tag[idx+1:]
		}
	}
	return field
}

func ruleMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "обязательное поле"
	case "min":
		return fmt.Sprintf("минимум %s", e.Param())
	case "max":
		return fmt.Sprintf("максимум %s", e.Param())
	case "email":
		return "неверный email"
	case "alphanum":
		return "только буквы и цифры"
	case "gte":
		return fmt.Sprintf("не меньше %s", e.Param())
	default:
		return fmt.Sprintf("не прошло проверку %s", e.Tag())
	}
}

var _ = reflect.TypeOf
