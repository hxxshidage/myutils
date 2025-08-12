package urtime

import "github.com/petermattis/goid"

func GoId() int64 {
	return goid.Get()
}
