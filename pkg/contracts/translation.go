package contracts

type Translator interface {
	TranslateWithData(messageKey string, templateData map[string]any) string
	Translate(messageKey string) string
}
