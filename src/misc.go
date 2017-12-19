package main

import (
	"strconv"
)

// From: https://stackoverflow.com/questions/19101419/go-golang-formatfloat-convert-float-number-to-string
func floatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
