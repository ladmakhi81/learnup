package validatorv10

import (
	"github.com/go-playground/validator/v10"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/types"
	"reflect"
)

type ValidatorSvc struct {
	core           *validator.Validate
	translationSvc translations.Translator
}

func NewValidatorSvc(
	core *validator.Validate,
	translationSvc translations.Translator,
) *ValidatorSvc {
	return &ValidatorSvc{
		core:           core,
		translationSvc: translationSvc,
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
		return svc.translationSvc.TranslateWithData(
			"common.errors.required_validation",
			map[string]any{
				"Name": fieldName,
			},
		)
	case "max":
		return svc.translationSvc.TranslateWithData(
			"common.errors.max_validation",
			map[string]any{
				"Name": fieldName,
				"Len":  err.Param(),
			},
		)
	case "min":
		return svc.translationSvc.TranslateWithData(
			"common.errors.min_validation",
			map[string]any{
				"Name": fieldName,
				"Len":  err.Param(),
			},
		)
	case "numeric":
		return svc.translationSvc.TranslateWithData(
			"common.errors.numeric_validation",
			map[string]any{
				"Name": fieldName,
			},
		)
	case "len":
		return svc.translationSvc.TranslateWithData(
			"common.errors.len_validation",
			map[string]any{
				"Name": fieldName,
				"Len":  err.Param(),
			},
		)
	default:
		return svc.translationSvc.TranslateWithData(
			"common.errors.unknown_validation",
			map[string]any{
				"Name": fieldName,
				"Tag":  err.Tag(),
			},
		)
	}
}
