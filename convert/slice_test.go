package uconv

import (
	"fmt"
	"testing"
)

func TestSliceConvert(t *testing.T) {
	src := []string{"1", "2", "3"}
	r := SliceConvert[string, int](src, S2i)
	fmt.Printf("%v", r)

	rr := SliceI2s(r)
	fmt.Printf("%v", rr)
}

func TestSliceConvertPost(t *testing.T) {
	src := []int{1, 2, 3}
	r := SliceConvertPost[int, string](src, I2s, func(i int, s string) {
		fmt.Printf("i:%d => s:%s\n", i, s)
	})
	fmt.Printf("%v", r)
}
