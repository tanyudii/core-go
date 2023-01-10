package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/tanyudii/core-go/errutil"
)

type Service interface {
	GetConfig() *Config
	GetValidate() *validator.Validate
	ValidateStruct(val interface{}) error
}

type service struct {
	cfg      *Config
	validate *validator.Validate
}

func NewValidator(args ...ConfigFunc) Service {
	v := &service{cfg: generateConfig(args...), validate: validator.New()}
	v.init()
	return v
}

func (s *service) GetConfig() *Config {
	return s.cfg
}

func (s *service) GetValidate() *validator.Validate {
	return s.validate
}

func (s *service) ValidateStruct(val interface{}) error {
	err := s.validate.Struct(val)
	if err == nil {
		return nil
	}
	fields := errutil.ErrorField{}
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			fields[s.transformField(e.StructNamespace())] = e.Translate(s.cfg.trans)
		}
	}
	return errutil.NewBadRequestErrorUsingFieldsOrNil(fields)
}
