package main

import (
	"math/rand"
)

func getRandomBit(len int) []byte {
	token := make([]byte, len)
	rand.Read(token)
	return token
}
