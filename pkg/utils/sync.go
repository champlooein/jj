package utils

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/golang/glog"
)

const defaultConcurrency = 10

// GoWaits 并发执行所有 fns 并等待执行结束.
func GoWaits(ctx context.Context, fns ...GoFunc) {
	GoWaitsWithLimit(defaultConcurrency, fns...)
}

// GoWaitsWithLimit 并发执行所有 fns 并等待执行结束, concurrency 表示允许同时进行的并发数. 小于等于 0 则使用默认值.
func GoWaitsWithLimit(concurrency int, fns ...GoFunc) {
	if concurrency <= 0 {
		concurrency = defaultConcurrency
	}

	pool := make(chan struct{}, concurrency) // 并发控制
	defer close(pool)

	wg := sync.WaitGroup{}
	wg.Add(len(fns))

	for _, fn := range fns {
		pool <- struct{}{}
		f := fn

		go func() {
			defer func() {
				<-pool
				wg.Done()
			}()

			f.goWithRecover()
		}()
	}

	wg.Wait()
}

// Go 启动一个goroutine, 包装recover和日志.
func Go(fn GoFunc) {
	go func() {
		fn.goWithRecover()
	}()
}

type GoFunc func()

func (fn GoFunc) goWithRecover() {
	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("Got panic: %+v\n\n%s", err, debug.Stack())
			return
		}
	}()

	if fn == nil {
		glog.Error("fn cannot be nil")
		return
	}

	fn()
}
