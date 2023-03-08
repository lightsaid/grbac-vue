package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 随机生成一个数，介于[min-max)
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
