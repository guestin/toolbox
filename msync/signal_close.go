package msync

import (
	"sync/atomic"
)

// 关闭信号
type SignalClose struct {
	closeFlag int32     //flag
	sigChan   chan bool //broadcast only
}

func NewSignalClose() *SignalClose {
	return &SignalClose{
		closeFlag: 0,
		sigChan:   make(chan bool),
	}
}

//判断信号是否被发送
func (this *SignalClose) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

//获取 signal chan ,该 chan 只在被关闭时返回数据
func (this *SignalClose) GetSignalChan() <-chan bool {
	return this.sigChan
}

// 通知所有关注close chan 的 goroutine, 成功执行时,返回 true, 如果已经被关闭,则返回 false
func (this *SignalClose) SendCloseSignal() bool {
	//fast path
	if !atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		return false
	}
	close(this.sigChan)
	return true
}
