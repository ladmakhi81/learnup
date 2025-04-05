package translations

type Translator interface {
	Translate(messageKey string, templateData map[string]any) string
}
