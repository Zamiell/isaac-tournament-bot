package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// From: https://stackoverflow.com/questions/19101419/go-golang-formatfloat-convert-float-number-to-string
func floatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// From: http://golangcookbook.blogspot.com/2012/11/generate-random-number-in-given-range.html
func getRandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	max++
	return rand.Intn(max-min) + min
}

func deleteFromSlice(a []string, i int) []string {
	return append(a[:i], a[i+1:]...)
}

func sliceToString(slice []string) string {
	return strings.Join(slice, ",")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func stringToSlice(str string) []string {
	// "strings.Split" will return a slice of length one if fed an empty string
	if str == "" {
		return make([]string, 0)
	}

	return strings.Split(str, ",")
}
