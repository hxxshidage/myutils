package uconv

import (
	utype "myutils/type"
)

type Converter[T any, R any] func(int, T) R

func SliceConvert[T any, R any](src []T, converter Converter[T, R]) []R {
	target := make([]R, len(src))
	for idx, t := range src {
		target[idx] = converter(idx, t)
	}

	return target
}

type PostFunc[T any, R any] func(T, R)

func SliceConvertPost[T any, R any](src []T, converter Converter[T, R], pf PostFunc[T, R]) []R {
	target := make([]R, len(src))
	for idx, t := range src {
		r := converter(idx, t)
		if nil != pf {
			pf(t, r)
		}

		target[idx] = r
	}

	return target
}

func Slice2s[T any](src []T, converter Converter[T, string]) []string {
	return SliceConvert(
		src,
		converter,
	)
}

func SliceI2s(src []int) []string {
	return SliceConvert(src, func(idx int, val int) string {
		return utype.I2s(val)
	})
}

func SliceS2i(src []string) []int {
	return SliceConvert(src, func(idx int, val string) int {
		return utype.S2i(val)
	})
}

// 开发者自行承担数据类型的正确性
func SliceMust2s(src any) []string {
	return SliceConvert[any, string](src.([]any), func(_ int, a any) string {
		return a.(string)
	})
}

type KeyMapper[T any, K comparable] func(T) K
type ValueMapper[T any, K comparable, V any] func(T, K) V

func Slice2map[T any, K comparable, V any](src []T, km KeyMapper[T, K], vm ValueMapper[T, K, V]) map[K]V {
	m := make(map[K]V, len(src))
	for _, t := range src {
		k := km(t)
		m[k] = vm(t, k)
	}

	return m
}

func Slice2itMap[T any, K comparable](src []T, km KeyMapper[T, K]) map[K]T {
	return Slice2map(
		src,
		km,
		func(t T, k K) T {
			return t
		},
	)
}

type Filter[T any] func(T) bool

func SliceFilter[T any](src []T, f Filter[T]) []T {
	if len(src) == 0 {
		return make([]T, 0)
	}

	newSli := make([]T, 0, len(src))

	for _, t := range src {
		if f(t) {
			newSli = append(newSli, t)
		}
	}

	return newSli
}

func SliceGroupBy[T any, K comparable, V any](src []T, km KeyMapper[T, K], vm ValueMapper[T, K, V]) map[K][]V {
	if len(src) == 0 {
		return make(map[K][]V)
	}

	m := make(map[K][]V)
	for _, t := range src {
		k := km(t)
		v := vm(t, k)

		if _, ok := m[k]; !ok {
			m[k] = []V{v}
		} else {
			m[k] = append(m[k], v)
		}
	}

	return m
}

func SliceGroupByIt[T any, K comparable](src []T, km KeyMapper[T, K]) map[K][]T {
	return SliceGroupBy(src, km, func(t T, k K) T {
		return t
	})
}
