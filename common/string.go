package common

import "strings"

func ParseStringToSliceBySeparator(str string, sep string) []string {
	if string(str) == "" {
		return nil
	}
	return strings.Split(str, sep)
}
