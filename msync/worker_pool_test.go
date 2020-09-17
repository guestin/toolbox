package msync

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func BenchmarkWorkerPool(b *testing.B) {
	caseFn := func() {
		cancel, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		pool := NewWorkerPool(cancel, "tester", runtime.NumCPU())
		pool.SetIdleMax(time.Second * 5)
		defer pool.Wait()
		for i := 0; i < 100; i++ {
			s := fmt.Sprintf("task:%d", i)
			pool.Submit(func() {
				println(s, " enter")
				println(s, " exit")
				time.Sleep(time.Second * 1)
			})
		}
	}
	b.ResetTimer()
	for i := 0; i < 20; i++ {
		caseFn()
		println("\n\n\n\n")
	}
}
