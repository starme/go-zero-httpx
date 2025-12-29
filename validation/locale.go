package validation

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/zh"
)

// TranslationRegistrationFunc registers translations for a validator and translator pair.
type TranslationRegistrationFunc func(*validator.Validate, ut.Translator) error

// DefaultTranslationRegistry maps locale codes to their registration helpers.
var DefaultTranslationRegistry = map[string]TranslationRegistrationFunc{
	"zh": func(v *validator.Validate, trans ut.Translator) error {
		return zh.RegisterDefaultTranslations(v, trans)
	},
	"en": func(v *validator.Validate, trans ut.Translator) error {
		return en.RegisterDefaultTranslations(v, trans)
	},
}

// RegisterLocaleTranslation adds a custom translation registration function for a locale.
func RegisterLocaleTranslation(lang string, fn TranslationRegistrationFunc) {
	if lang == "" || fn == nil {
		return
	}

	DefaultTranslationRegistry[lang] = fn
}

// RegisterTranslationsForLang runs the registered translation function for the specified locale.
func RegisterTranslationsForLang(v *validator.Validate, trans ut.Translator, lang string) error {
	if fn, ok := DefaultTranslationRegistry[lang]; ok && fn != nil {
		return fn(v, trans)
	}

	return fmt.Errorf("no translation registered for language %q", lang)
}
