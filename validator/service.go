package validator

import (
	"context"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/tanyudii/core-go/ectxutil"
	"github.com/tanyudii/core-go/errutil"
)

type Service interface {
	ValidateStruct(val interface{}) error
	ValidateStructWithContext(ctx context.Context, val interface{}) error
	ValidateStructWithLang(lang string, val interface{}) error
}

type service struct {
	cfg      *Config
	validate *validator.Validate
	uni      *ut.UniversalTranslator
}

func NewValidator(args ...ConfigFunc) Service {
	v := &service{cfg: generateConfig(args...)}
	v.init()
	return v
}

func (s *service) ValidateStruct(val interface{}) error {
	return s.ValidateStructWithLang(DefaultLocaleName, val)
}

func (s *service) ValidateStructWithContext(ctx context.Context, val interface{}) error {
	lang := ectxutil.GetAcceptLanguage(ctx, DefaultLocaleName)
	return s.ValidateStructWithLang(lang, val)
}

func (s *service) ValidateStructWithLang(lang string, val interface{}) error {
	err := s.validate.Struct(val)
	if err == nil {
		return nil
	}
	translatorFn := s.getTranslator(lang)
	fields := errutil.ErrorField{}
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			fields[s.transformField(e.StructNamespace())] = e.Translate(translatorFn)
		}
	}

	return errutil.NewBadRequestErrorUsingFieldsOrNil(fields)
}
