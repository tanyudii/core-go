package common

import (
	"fmt"
	"regexp"
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

func IsStringInSeparatedString(val, check, sep string) bool {
	return SeparatedStringToMapBool(val, sep)[check]
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func AppendSeparatedStrings(current, sep string, values ...string) string {
	newValue := strings.Join(values, sep)
	if current == "" {
		return newValue
	}
	return current + sep + newValue
}
