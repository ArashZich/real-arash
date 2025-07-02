package codegen

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomNumber() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	code := rand.Intn(max-min+1) + min
	c := strconv.Itoa(code)
	return c
}
