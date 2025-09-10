package internal

import (
	"strings"
	"time"
)

type multiPartString string

func FormatDateYYYYMMDD(t time.Time) string {
	return t.Format("2006/01/02")
}

func (originalString multiPartString) containsAny(substrings []string) bool {
	for _, subString := range substrings {
		if strings.Contains(string(originalString), subString) {
			return true
		}
	}
	return false
}
