package msync

import (
	"github.com/guestin/mob/merrors"
	"sync"
	"sync/atomic"
)

// best practice for wait_group
type AsyncTaskGroup struct {
	wg      *sync.WaitGroup
	taskCnt int32
}

func NewAsyncTaskGroup() *AsyncTaskGroup {
	return NewAsyncTaskGroup2(new(sync.WaitGroup))
}

func NewAsyncTaskGroup2(wg *sync.WaitGroup) *AsyncTaskGroup {
	merrors.Assert(wg != nil, "waitGroup is nil")
	return &AsyncTaskGroup{
		wg:      wg,
		taskCnt: 0,
	}
}

func (this *AsyncTaskGroup) AddTask(exe func()) {
	merrors.Assert(exe != nil, "exe cant be nil")
	atomic.AddInt32(&this.taskCnt, 1)
	this.wg.Add(1)
	go func() {
		defer func() {
			this.wg.Done()
			atomic.AddInt32(&this.taskCnt, -1)
		}()
		exe()
	}()
}

func (this *AsyncTaskGroup) GetTotalTask() int {
	return int(atomic.LoadInt32(&this.taskCnt))
}

func (this *AsyncTaskGroup) Wait() {
	this.wg.Wait()
}
