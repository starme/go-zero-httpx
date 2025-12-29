package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// TransFn defines a translation helper compatible with validator.TranslationFunc.
type TransFn validator.TranslationFunc

// RegisFn builds a validator translation registration function.
type RegisFn func(tag string, translation string, override bool) validator.RegisterTranslationsFunc

// CustomValidator holds metadata for a custom validator and its translations.
type CustomValidator struct {
	tag                  string
	validateFunc         validator.Func
	translation          string
	override             bool
	customRegistrationFn validator.RegisterTranslationsFunc
	customTranslateFn    validator.TranslationFunc
}

// NewCustomValidator creates an empty CustomValidator ready for chained configuration.
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{}
}

// Validator extends the go-playground validator and stores translation helpers.
type Validator struct {
	*validator.Validate

	trans                ut.Translator
	customValidator      []*CustomValidator
	customTagNameFn      validator.TagNameFunc
	customTranslateFn    TransFn
	customRegistrationFn RegisFn
}

var (
	validatorOnce    sync.Once
	defaultValidator *Validator
)

// NewValidator returns a shared Validator configured with defaults and optional translator.
func NewValidator(trans ut.Translator) *Validator {
	validatorOnce.Do(func() {
		defaultValidator = &Validator{
			Validate:             validator.New(),
			trans:                trans,
			customValidator:      []*CustomValidator{},
			customTagNameFn:      defaultTagNameFunc,
			customTranslateFn:    defaultTranslateFunc,
			customRegistrationFn: defaultRegistrationFunc,
		}
	})

	return defaultValidator
}

// AddCustom registers custom validation and translation rules on the validator.
func (v *Validator) AddCustom(customValidator ...*CustomValidator) (*Validator, error) {
	var errs []error
	for _, t := range customValidator {
		if err := v.RegisterValidation(t.tag, t.validateFunc); err != nil {
			continue
		}

		if t.customRegistrationFn == nil {
			t.customRegistrationFn = v.customRegistrationFn(t.tag, t.translation, t.override)
		}

		if t.customTranslateFn == nil {
			t.customTranslateFn = validator.TranslationFunc(v.customTranslateFn)
		}

		if err := v.RegisterTranslation(t.tag, v.trans, t.customRegistrationFn, t.customTranslateFn); err != nil {
			errs = append(errs, fmt.Errorf("register translation for tag %q failed: %w", t.tag, err))
		}
	}

	if len(errs) > 0 {
		return v, errors.Join(errs...)
	}

	v.customValidator = append(v.customValidator, customValidator...)

	return v, nil
}

// defaultRegistrationFunc registers the tag translation when no override is needed.
func defaultRegistrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		return ut.Add(tag, translation, override)
	}
}

// defaultTranslateFunc returns the translated validation message or falls back to the error string.
func defaultTranslateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		fmt.Printf("警告: 翻译字段错误: %#v\n", fe)
		return fe.(error).Error()
	}

	return t
}

// defaultTagNameFunc uses the json tag as the field name for validation errors.
func defaultTagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	if name == "-" {
		return ""
	}

	return name
}

// NewTranslator builds a universal translator with the requested default language.
func NewTranslator(defaultLang string, supportLocales ...locales.Translator) ut.Translator {
	translator := ut.New(supportLocales[0], supportLocales[1:]...)

	trans, found := translator.GetTranslator(defaultLang)
	if !found {
		fmt.Println("translator not found")
	}

	return trans
}
