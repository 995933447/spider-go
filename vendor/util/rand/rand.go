package rand

import (
	"math/rand"
	"time"
)

func RandInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}

func RandIntN(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}