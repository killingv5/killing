package utils

import (
	"math/rand"
)

const RANDOM_TRY_MULTIPLE = 2

func GenRandomInt(uplimt int) int {
	//rand.Seed(time.Now().UnixNano())
	index := rand.Intn(uplimt)
	return index
}
