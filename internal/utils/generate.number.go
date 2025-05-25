package utils

import (
	"math/rand"
)

func GenerateNumber() (string, error) {
	char := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	code := make([]byte, 5)
	for i := range code {
		code[i] = char[rand.Intn(len(char))]
	}

	return string(code), nil
}
