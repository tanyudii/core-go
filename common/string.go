package common

import (
	"fmt"
	"strings"
)

func ParseStringToSliceBySeparator(str string, sep string) []string {
	if string(str) == "" {
		return nil
	}
	return strings.Split(str, sep)
}

func RFC5322Format(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

func SeparatedStringToMapBool(val, sep string) map[string]bool {
	if val == "" {
		return nil
	}
	ret := make(map[string]bool)
	for _, v := range strings.Split(val, sep) {
		ret[v] = true
	}
	return ret
}

func AppendSeparatedString(current, new, sep string) string {
	if current == "" {
		return new
	}
	return current + sep + new
}

func IsStringInSeparatedString(val, check, sep string) bool {
	return SeparatedStringToMapBool(val, sep)[check]
}
