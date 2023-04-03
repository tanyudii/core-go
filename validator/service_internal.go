package validator

import (
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"reflect"
	"regexp"
	"strings"
)

func (s *service) init() {
	s.initUni()

	s.validate = validator.New()
	s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("label")
	})

	for name, fn := range s.cfg.mapRegisterTranslation {
		translator, _ := s.uni.GetTranslator(name)
		_ = fn(s.validate, translator)
	}

}

func (s *service) initUni() {
	if s.uni != nil {
		return
	}
	var translations []locales.Translator
	for _, trans := range s.cfg.mapLocalesTranslator {
		translations = append(translations, trans)
	}
	if len(translations) > 1 {
		s.uni = ut.New(translations[0], translations[1:]...)
	} else {
		s.uni = ut.New(translations[0])
	}
}

func (s *service) getTranslator(localeName string) ut.Translator {
	transFn, ok := s.uni.GetTranslator(localeName)
	if !ok {
		transFn, _ = s.uni.GetTranslator(DefaultLocaleName)
	}
	return transFn
}

func (s *service) transformField(field string) string {
	rgx, _ := regexp.Compile("\\[(.*?)]")
	fields := strings.Split(field, ".")
	var attrs []string
	for _, attr := range fields[1:] {
		match := rgx.FindStringSubmatch(attr)
		if len(match) == 2 {
			attrs = append(attrs, strcase.ToLowerCamel(rgx.ReplaceAllString(attr, "")), match[1])
		} else {
			attrs = append(attrs, strcase.ToLowerCamel(attr))
		}
	}
	return strings.Join(attrs, ".")
}
