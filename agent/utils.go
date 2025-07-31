package main

import (
	"math/rand"
	"time"
)

func tickAmount(min, max int) time.Duration {
	if min == max {
		return time.Millisecond * time.Duration(min)
	}
	return time.Millisecond * time.Duration(rand.Intn(max-min)+max)
}
