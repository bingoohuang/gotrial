package main

import (
	"github.com/averagesecurityguy/random"
)

// RandomFloat64 returns a random float64
func RandomFloat64() float64 {
	i, _ := random.Int64()
	return float64(i) * 1.0842021724855043e-19
}

// RandomInt returns a random int
func RandomInt(n uint64) int {
	i, _ := random.Uint64Range(0, n)
	return int(i)
}

// RandomInt64 returns a random int64
func RandomInt64() int64 {
	i, _ := random.Int64()
	return i
}

// RandomString returns a random string
func RandomString(n uint64) string {
	s, _ := random.AlphaNum(n)
	return s
}
