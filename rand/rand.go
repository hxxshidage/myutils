package urand

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	words = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	wLen  = len(words)
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

//var src = rand.NewSource(time.Now().UnixNano())

var srcPool = sync.Pool{
	New: func() any {
		// 源码中明确注释, 非并发安全
		// NewSource returns a new pseudo-random [Source] seeded with the given value.
		// Unlike the default [Source] used by top-level functions, this source is not
		// safe for concurrent use by multiple goroutines.
		// The returned [Source] implements [Source64].
		return rand.NewSource(time.Now().UnixNano())
	},
}

func RandStr(len uint16) string {
	appender := strings.Builder{}
	rLen := int(len)
	appender.Grow(rLen)

	src := srcPool.Get().(rand.Source)
	defer srcPool.Put(src)

	for i, cache, remain := rLen-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < wLen {
			appender.WriteByte(words[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return appender.String()
}

func RandInt(start, end int) int {
	return rand.Intn(end-start+1) + start
}
