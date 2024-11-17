package tool

import (
	"strconv"
	"strings"
)

func ConvertToInArray(text string) []int {
	var intSlice []int
	for _, s := range strings.Split(strings.Trim(text, "[]"), ",") {
		if num, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			intSlice = append(intSlice, num)
		}
	}
	return intSlice
}
