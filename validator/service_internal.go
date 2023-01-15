package validator

import (
	"github.com/iancoleman/strcase"
	"reflect"
	"regexp"
	"strings"
)

func (s *service) init() {
	s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("label")
	})
	_ = s.cfg.registerTransFunc(s.validate, s.cfg.trans)
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
