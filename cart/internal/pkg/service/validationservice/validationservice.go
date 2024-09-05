package validationservice

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
)

var validate = validator.New()

func Validate(v interface{}) error {
	err := validate.Struct(v)
	resultErrorMessage := ""
	if err != nil {
		// Если есть ошибки валидации, они будут возвращены здесь
		for _, err := range err.(validator.ValidationErrors) {
			resultErrorMessage += fmt.Sprintf("Error: Field '%s' failed validation with tag '%s'\n", err.Field(), err.Tag())
		}
		log.Printf("Errors found during validation: \n%s", resultErrorMessage)
		return errors.New(resultErrorMessage)
	}
	return nil
}
