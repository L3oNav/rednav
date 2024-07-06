package utils

import (
	"math/rand"
	"time"
)

func FlattenList(inputList [][]byte) []byte {
	var flatList []byte
	for _, sublist := range inputList {
		flatList = append(flatList, sublist...)
	}
	return flatList
}

func GenerateAlphanumericString() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 40)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
