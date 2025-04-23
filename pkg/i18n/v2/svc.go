package i18nv2

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path"
)

type I18nTranslatorSvc struct {
	localizer *i18n.Localizer
}

func setupLocalizer() (*i18n.Localizer, error) {
	locales := map[string]struct {
		langTag  language.Tag
		langText string
	}{
		"fa": {
			langTag:  language.Persian,
			langText: "fa",
		},
		"en": {
			langTag:  language.English,
			langText: "en",
		},
	}
	defaultLocale := "fa"
	localizationBundle := i18n.NewBundle(locales[defaultLocale].langTag)
	localizationBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	rootDir, rootDirErr := os.Getwd()
	if rootDirErr != nil {
		return nil, rootDirErr
	}
	translationFolderPath := path.Join(rootDir, "translations")
	localizationBundle.MustLoadMessageFile(path.Join(translationFolderPath, "fa.json"))
	localizationBundle.MustLoadMessageFile(path.Join(translationFolderPath, "en.json"))
	return i18n.NewLocalizer(localizationBundle, locales[defaultLocale].langText), nil
}

func NewI18nTranslatorSvc() (*I18nTranslatorSvc, error) {
	localizer, err := setupLocalizer()
	if err != nil {
		return nil, err
	}
	return &I18nTranslatorSvc{
		localizer: localizer,
	}, nil
}

func (svc I18nTranslatorSvc) TranslateWithData(messageKey string, templateData map[string]any) string {
	return svc.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageKey,
		TemplateData: templateData,
	})
}

func (svc I18nTranslatorSvc) Translate(messageKey string) string {
	return svc.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageKey,
	})
}
