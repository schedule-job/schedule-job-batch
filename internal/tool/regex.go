package tool

import (
	"regexp"
)

func FindWords(text string, reg string) []string {
	regex, _ := regexp.Compile(reg)
	data := []string{}
	for _, fn := range regex.FindAllStringSubmatch(text, -1) {
		for i, f := range fn {
			if i == 0 {
				continue
			}
			data = append(data, f)
		}

	}
	return data
}
