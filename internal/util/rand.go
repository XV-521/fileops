package util

import (
	"math/rand"
	"strconv"
)

func GetRand(digit int) string {
	result := ""
	for range digit {
		num := rand.Intn(10)
		result += strconv.Itoa(num)
	}
	return result
}
