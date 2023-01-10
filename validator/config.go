package validator

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

const (
	DefaultLocaleName = "en"
)

var (
	DefaultLocale = en.New()
)

type Config struct {
	locale            locales.Translator
	localeName        string
	uni               *ut.UniversalTranslator
	trans             ut.Translator
	registerTransFunc TransFunc
}

type ConfigFunc func(c *Config)

type TransFunc func(v *validator.Validate, trans ut.Translator) (err error)

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{locale: DefaultLocale, localeName: DefaultLocaleName}
	for i := range args {
		args[i](c)
	}
	if c.uni == nil {
		c.uni = c.defaultUni()
	}
	if c.trans == nil {
		c.trans = c.defaultTrans()
	}
	if c.registerTransFunc == nil {
		c.registerTransFunc = entranslations.RegisterDefaultTranslations
	}
	return c
}

func Locale(l locales.Translator, name string) ConfigFunc {
	return func(c *Config) {
		c.locale = l
		c.localeName = name
	}
}

func Uni(uni *ut.UniversalTranslator) ConfigFunc {
	return func(c *Config) {
		c.uni = uni
	}
}

func Trans(t ut.Translator) ConfigFunc {
	return func(c *Config) {
		c.trans = t
	}
}

func RegisterTransFunc(fn TransFunc) ConfigFunc {
	return func(c *Config) {
		c.registerTransFunc = fn
	}
}
