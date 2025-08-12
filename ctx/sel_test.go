package uctx

import (
	"fmt"
	"testing"
)

func TestSel(t *testing.T) {
	fmt.Println("start")

	// select一直阻塞, fatal error: all goroutines are asleep - deadlock!
	select {}

	fmt.Println("end...")
}
