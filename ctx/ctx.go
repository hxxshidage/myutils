package uctx

import (
	"context"
	"github.com/pkg/errors"
	"reflect"
	"runtime"
	"time"
)

type Run[T any] func() (T, error)

func RunWithTimeout[T any](ctx context.Context, timeout time.Duration, r Run[T]) (T, error) {
	if nil == ctx {
		ctx = context.Background()
	}

	cclCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rch := make(chan struct {
		val T
		err error
	}, 1)

	go func() {
		defer func() {
			if rv := recover(); rv != nil {
				rch <- struct {
					val T
					err error
				}{err: errors.Errorf("%v", rv)}
			}
		}()

		v, e := r()

		rch <- struct {
			val T
			err error
		}{val: v, err: e}
	}()

	select {
	case res := <-rch:
		return res.val, res.err
	case <-cclCtx.Done():
		var zero T
		if ce := cclCtx.Err(); ce != nil {
			rName := runtime.FuncForPC(reflect.ValueOf(r).Pointer()).Name()
			return zero, errors.Wrapf(ce, "execute runnable:%s failed", rName)
		}

		return zero, nil
	}
}
