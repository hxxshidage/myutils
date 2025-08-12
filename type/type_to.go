package utype

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

func MustStr(val any) string {
	return val.(string)
}

func MustIntF64(val any) int {
	return int(val.(float64))
}

func MustInt(val any) int {
	return val.(int)
}

func MustStrMap(val any) map[string]any {
	return val.(map[string]any)
}

// str => int
func S2i(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	} else {
		panic(errors.Wrap(err, fmt.Sprintf("parse str num:%s to int32 failed", s)))
	}
}

// str => int64
func S2i64(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	} else {
		panic(errors.Wrap(err, fmt.Sprintf("parse str num:%s to int64 failed", s)))
	}
}

// int => str
func I2s(i int) string {
	return strconv.Itoa(i)
}

// int64 => str
func I642s(i int64) string {
	return strconv.FormatInt(i, 10)
}

func S2b(s string) bool {
	if b, err := strconv.ParseBool(s); err != nil {
		return false
	} else {
		return b
	}
}

func B2s(b bool) string {
	return strconv.FormatBool(b)
}

func B2i(b bool) int {
	i := 0
	if b {
		i = 1
	}

	return i
}
