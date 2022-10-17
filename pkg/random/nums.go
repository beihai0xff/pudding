package random

import (
	"math/rand"
	"time"
)

func GetRand(start, end int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(end-start) + start
}
