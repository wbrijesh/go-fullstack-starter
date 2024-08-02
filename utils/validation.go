package utils

import "strings"

func IsValidEmail(email string) bool {
	if !strings.Contains(email, "@") {
		// must have an @ symbol
		return false
	} else if !strings.Contains(email, ".") {
		// must have a period
		return false
	} else if strings.Index(email, "@") > strings.Index(email, ".") {
		// should come before the period
		return false
	} else if strings.Index(email, "@") == 0 {
		// should not be the first character
		return false
	} else if strings.Index(email, ".") == len(email)-1 {
		// period should not be the last character
		return false
	} else if strings.Index(email, ".")-strings.Index(email, "@") == 1 {
		// period should not come directly after @
		return false
	} else {
		// all checks passed
		return true
	}
}
