package models

import (
	"strings"
)

/*
	Miscellaneous subroutines
*/

func stringToSlice(a string) []string {
	slice := strings.Split(a, ",")
	if a == "" {
		slice = make([]string, 0)
	}
	return slice
}

func sliceToString(slice []string) string {
	return strings.Join(slice, ",")
}
