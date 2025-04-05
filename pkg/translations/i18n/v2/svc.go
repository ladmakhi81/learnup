package i18nv2

import "github.com/nicksnyder/go-i18n/v2/i18n"

type I18nTranslatorSvc struct {
	localizer *i18n.Localizer
}

func NewI18nTranslatorSvc(localizer *i18n.Localizer) *I18nTranslatorSvc {
	return &I18nTranslatorSvc{
		localizer: localizer,
	}
}

func (svc I18nTranslatorSvc) Translate(messageKey string, templateData map[string]any) string {
	return svc.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageKey,
		TemplateData: templateData,
	})
}
