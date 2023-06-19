package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func deleteFromSlice(a []string, i int) []string {
	return append(a[:i], a[i+1:]...)
}

// From: https://stackoverflow.com/questions/19101419/go-golang-formatfloat-convert-float-number-to-string
func floatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// Returns the random array element and the randomly chosen index.
func getRandomArrayElement[T any](array []T) (T, int) {
	if len(array) == 0 {
		log.Panic("Failed to get a random array element since the provided array was empty.")
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(array))
	return array[index], index
}

// Returns a random integer between min and max. It is inclusive on both ends. For example,
// `getRandomInt(1, 3)` will return 1, 2, or 3.
func getRandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	max++
	return rand.Intn(max-min) + min
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
	// "strings.Split" will return a slice of length one if fed an empty string.
	if str == "" {
		return make([]string, 0)
	}

	return strings.Split(str, ",")
}
