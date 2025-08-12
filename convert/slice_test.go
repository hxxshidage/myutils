package uconv

import (
	"fmt"
	utype "github.com/hxxshidage/myutils/type"
	"testing"
)

func TestSliceConvert(t *testing.T) {
	src := []string{"1", "2", "3"}
	r := SliceConvert[string, int](src, func(idx int, item string) int {
		return utype.S2i(item)
	})
	fmt.Printf("%v", r)

	rr := SliceI2s(r)
	fmt.Printf("%v", rr)
}

func TestSliceConvertPost(t *testing.T) {
	src := []int{1, 2, 3}
	r := SliceConvertPost[int, string](src, func(idx int, item int) string {
		return utype.I2s(item)
	}, func(in int, out string) {
		fmt.Printf("i:%d => s:%s\n", in, out)
	})
	fmt.Printf("%v", r)
}
