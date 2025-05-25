package utils

import "strings"

func LowerCaseStrings(input []string) []string {
	output := make([]string, len(input))
	for i, s := range input {
		output[i] = strings.ToLower(s)
	}
	return output
}
