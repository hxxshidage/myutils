package uctx

import (
	"context"
	"github.com/labstack/gommon/log"
	"testing"
	"time"
)

func TestCtxTimeout(t *testing.T) {
	log.Info("start")

	ttCtx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	ch := make(chan struct{}, 1)

	go func() {
		time.AfterFunc(time.Second*3, func() {
			log.Info("write data to ch")
			ch <- struct{}{}
		})
	}()

	select {
	case <-ch:
		log.Info("read data from ch")
	case <-ttCtx.Done():
		log.Warnf("timeout, reason:%v", ttCtx.Err())
	}

	log.Info("end...")
}

/*
type Context interface {
	// 如果上下文设置了截止时间（比如通过 context.WithDeadline 或 context.WithTimeout 创建），那么返回的 ok 为 true，并且 deadline 表示这个上下文将在该时间点被取消
	// 如果上下文没有设置截止时间（比如是 context.Background() 或 context.TODO()），则 ok 为 false，deadline 是零值（即 time.Time{}）
	// 作用:在各个子context操作时, 可以知道还剩余多少时间可以操作, 以决定到底操不操作
	Deadline() (deadline time.Time, ok bool)

	// 当被主动取消或超时 "完成"时, 会关闭相应channel, 通过Done()函数监听这个只读channel(监听被关闭的通道, 会立刻返回零值)
	// 通常是监听ctx.Done(), 再根据ctx.Err()来判断原因
	Done() <-chan struct{}

	// 返回被取消的原因
	Err() error

	// 传值了
	Value(key any) any
}
*/

func TestRunWithTimeout(t *testing.T) {
	rst, err := RunWithTimeout(nil, time.Second*1, func() (string, error) {
		time.Sleep(2 * time.Second)
		return "", nil
	})

	if err != nil {
		log.Error(err)
	}

	log.Info(rst)

}
