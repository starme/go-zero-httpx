package validation

import (
	"testing"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRegisterTranslationsForLang_Default(t *testing.T) {
	cases := []struct {
		name   string
		locale locales.Translator
		lang   string
	}{
		{name: "zh", locale: zh.New(), lang: "zh"},
		{name: "en", locale: en.New(), lang: "en"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			trans := newTranslator(t, tt.locale)
			v := validator.New()

			require.NoError(t, RegisterTranslationsForLang(v, trans, tt.lang))
		})
	}
}

func TestRegisterTranslationsForLang_MissingLang(t *testing.T) {
	trans := newTranslator(t, zh.New())
	v := validator.New()

	err := RegisterTranslationsForLang(v, trans, "en-US")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no translation registered")
}

func TestRegisterLocaleTranslation_Dynamic(t *testing.T) {
	lang := "mock-lang"
	called := false
	RegisterLocaleTranslation(lang, func(v *validator.Validate, trans ut.Translator) error {
		called = true
		return nil
	})
	defer delete(DefaultTranslationRegistry, lang)

	trans := newTranslator(t, zh.New())
	v := validator.New()

	require.NoError(t, RegisterTranslationsForLang(v, trans, lang))
	require.True(t, called, "dynamic translation func should be invoked")
}

func newTranslator(t *testing.T, locale locales.Translator) ut.Translator {
	t.Helper()

	translator := ut.New(locale, locale)
	trans, found := translator.GetTranslator(locale.Locale())
	require.True(t, found, "%s translator should be available", locale.Locale())

	return trans
}
