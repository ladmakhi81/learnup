package validatorv10

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/ladmakhi81/learnup/types"
	"reflect"
)

type ValidatorSvc struct {
	core *validator.Validate
}

func NewValidatorSvc(core *validator.Validate) *ValidatorSvc {
	return &ValidatorSvc{
		core: core,
	}
}

func (svc *ValidatorSvc) Validate(dto any) *types.ClientError {
	err := svc.core.Struct(dto)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			jsonTagName := svc.getJsonTagName(dto, err.StructField())
			validationErrors = append(validationErrors, svc.generateErrMessage(jsonTagName, err))
		}
		return types.NewBadRequestDTOError(validationErrors)
	}
	return nil
}

func (svc *ValidatorSvc) getJsonTagName(dto any, fieldName string) string {
	t := reflect.TypeOf(dto)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	field, _ := t.FieldByName(fieldName)
	return field.Tag.Get("json")
}

func (svc *ValidatorSvc) generateErrMessage(fieldName string, err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", fieldName, err.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", fieldName, err.Param())
	case "numeric":
		return fmt.Sprintf("%s is not a valid number", fieldName)
	case "len":
		return fmt.Sprintf("%s must have length of %s", fieldName, err.Param())
	default:
		return fmt.Sprintf("%s with tag %s is not valid", fieldName, err.Tag())
	}
}
