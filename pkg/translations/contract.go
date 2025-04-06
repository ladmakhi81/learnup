package translations

type Translator interface {
	TranslateWithData(messageKey string, templateData map[string]any) string
	Translate(messageKey string) string
}
