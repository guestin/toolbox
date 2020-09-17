package msync

import (
	"container/list"
	"context"
	"fmt"
	"github.com/guestin/mob/merrors"
	"runtime"
	"sync"
	"time"
)

type WorkerFunc = func()

type WorkerPool struct {
	ctx         context.Context
	name        string
	maxWorkerN  int
	workQ       chan WorkerFunc
	startFlag   int32
	idle        time.Duration
	group       *AsyncTaskGroup
	taskBacklog *list.List
	cond        *sync.Cond
}

func NewWorkerPool(ctx context.Context, name string, expectWorkerNum int) *WorkerPool {
	if expectWorkerNum <= 0 {
		expectWorkerNum = runtime.NumCPU()
	}
	return &WorkerPool{
		ctx:         ctx,
		name:        name,
		maxWorkerN:  expectWorkerNum,
		workQ:       make(chan WorkerFunc, expectWorkerNum),
		startFlag:   0,
		idle:        time.Minute,
		group:       NewAsyncTaskGroup(),
		taskBacklog: list.New(),
		cond:        sync.NewCond(new(sync.Mutex)),
	}
}

func (this *WorkerPool) Wait() {
	this.group.Wait()
}

func (this *WorkerPool) SetIdleMax(idleTime time.Duration) {
	if idleTime <= time.Duration(0) {
		idleTime = time.Minute
	}
	this.idle = idleTime
}

func (this *WorkerPool) workerTask() {
	idleTimer := time.NewTimer(this.idle)
	defer idleTimer.Stop()
	//
	stealPipe := make(chan WorkerFunc)
	ctx, cancel := context.WithCancel(this.ctx)
	wg := sync.WaitGroup{}
	defer func() {
		cancel()
		this.cond.Broadcast()
		wg.Wait()
		close(stealPipe)
	}()
	wg.Add(1)
	go this.startWorkerStealer(ctx, stealPipe, &wg)
	//
	for {
		select {
		case workF, ok := <-stealPipe:
			if !ok || workF == nil {
				return
			}
			workF()
			idleTimer.Reset(this.idle)
		case workF, ok := <-this.workQ:
			if !ok || workF == nil {
				return
			}
			workF()
			idleTimer.Reset(this.idle)
		case <-idleTimer.C:
			return
		case <-this.ctx.Done():
			return
		}
	}
}

func (this *WorkerPool) startWorkerStealer(
	ctx context.Context,
	pipe chan<- WorkerFunc,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			{
				this.cond.L.Lock()
				if this.taskBacklog.Len() == 0 {
					this.cond.Wait()
					select {
					case <-ctx.Done():
						this.cond.L.Unlock()
						return
					default:
					}
					if this.taskBacklog.Len() == 0 {
						this.cond.L.Unlock()
						continue
					}
				}
				f := this.taskBacklog.Remove(this.taskBacklog.Front()).(WorkerFunc)
				this.cond.L.Unlock()
				pipe <- f
			}
		}
	}
}

func (this *WorkerPool) Submit(f WorkerFunc) {
	merrors.Assert(f != nil, fmt.Sprintf("worker pool [%s],submit nil task", this.name))
	select {
	case this.workQ <- f:
		if this.group.GetTotalTask() < this.maxWorkerN {
			this.group.AddTask(this.workerTask)
		}
	default:
		this.cond.L.Lock()
		defer this.cond.L.Unlock()
		this.taskBacklog.PushBack(f)
		this.cond.Broadcast()
	}
}

func (this *WorkerPool) GetName() string {
	return this.name
}
