package internal

import (
	"math/rand"
	"strconv"
	"strings"
)

func IsThisExt(filename string, ext string) bool {
	result := strings.Split(filename, ".")
	if len(result) < 2 {
		return ext == ""
	}
	return strings.ToLower(result[len(result)-1]) == strings.ToLower(ext)
}

func GetBasenameAndExt(filename string) (baseName string, ext string) {
	result := strings.Split(filename, ".")
	length := len(result)
	if length < 2 {
		return filename, ""
	}
	return strings.Join(result[0:length-1], "."), result[length-1]
}

func GetRand(digit int) string {
	result := ""
	for range digit {
		num := rand.Intn(10)
		result += strconv.Itoa(num)
	}
	return result
}
