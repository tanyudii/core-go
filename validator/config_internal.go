package validator

import ut "github.com/go-playground/universal-translator"

func (c *Config) defaultUni() *ut.UniversalTranslator {
	uni := ut.New(c.locale, c.locale)
	return uni
}

func (c *Config) defaultTrans() ut.Translator {
	trans, _ := c.uni.GetTranslator(c.localeName)
	return trans
}
