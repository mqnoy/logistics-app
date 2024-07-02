package cvalidator

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validator *validator.Validate
	once      sync.Once
)

const (
	ErrorValidator = "ERROR_VALIDATOR"
)

type ValidationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field: %s, code: %s, message: %s", e.Field, e.Code, e.Message)
}

func wrapError(field string, code string, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    code,
		Message: message,
	}
}

func init() {
	once.Do(func() {
		Validator = validator.New()

		// Register your custom validator function here
		Validator.RegisterValidation("normalize", normalizeString)
	})
}

func normalizeString(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	strLower := strings.ToLower(str)
	fl.Field().SetString(strLower)
	return true
}

func ValidateStruct(data interface{}) []ValidationError {
	var errValidator []ValidationError

	enLocale := en.New()
	uni := ut.New(enLocale)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(Validator, trans)

	err := Validator.Struct(data)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fieldName := e.StructNamespace()
			code := e.ActualTag()
			errorMessage := e.Translate(trans)
			errValidator = append(errValidator, wrapError(fieldName, code, errorMessage))
		}
	}

	return errValidator
}
