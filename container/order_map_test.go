package uconer

import (
	"fmt"
	"testing"
)

func TestNewOrderedMap(t *testing.T) {
	om := NewOrderedMapWithCap(8)

	setted := om.SetIfAbsent("1", "11")

	fmt.Printf("key:1, setted:%t\n", setted)

	setted = om.SetIfAbsent("1", "11")

	fmt.Printf("key:1, setted:%t\n", setted)

	om.Set("c", "cc")
	om.Set("b", "bb")
	om.Set("a", "aa")
	om.Set("e", "ee")

	if _, ok := om.Get("f"); !ok {
		om.Set("f", "ff")
	}

	computed := om.ComputeIfAbsent("z", func(key string) any {
		return "zz"
	})

	fmt.Printf("key:z, computed:%t\n", computed)

	computed = om.ComputeIfAbsent("z", func(key string) any {
		return "zz"
	})
	fmt.Printf("key:z, computed:%t\n", computed)

	om.Range(func(idx int, key string, val any) bool {
		fmt.Printf("key:%s, val:%v\n", key, val)
		return true
	})
}
