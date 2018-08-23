package main

import (
	"testing"
	"fmt"
)

func TestWallet(t *testing.T) {
	wallet := NewWallet()
	fmt.Println(wallet)
}
