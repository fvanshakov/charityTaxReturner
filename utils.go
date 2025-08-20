package main

import "strings"

type multiPartString string

func (originalString multiPartString) containsAny(substrings []string) bool {
	for _, subString := range substrings {
		if strings.Contains(string(originalString), subString) {
			return true
		}
	}
	return false
}
