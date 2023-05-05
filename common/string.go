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
