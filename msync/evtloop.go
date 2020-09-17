package msync

import (
	"container/list"
	"sync"
)

type EventHandlerFunc func(interface{})

var DiscardEventHandler = func(interface{}) {}

//event loop
type EvtLoop struct {
	*SignalClose
	evtList    *list.List
	evtHandler EventHandlerFunc
	cond       *sync.Cond
	name       string
}

//创建新的事件循环
func NewEvtLoop(name string) *EvtLoop {
	return &EvtLoop{
		evtList:     list.New(),
		evtHandler:  DiscardEventHandler,
		SignalClose: NewSignalClose(),
		cond:        sync.NewCond(&sync.Mutex{}),
		name:        name,
	}
}

func (this *EvtLoop) SetEventHandler(ie EventHandlerFunc) {
	if ie == nil {
		panic("EventHandlerFunc = nil")
	}
	this.evtHandler = ie
}

//一旦 StopLoop 返回,EventHandler将不会收到任何消息
//剩余的消息,全部会被丢弃
func (this *EvtLoop) StopLoop() bool {
	if !this.SendCloseSignal() {
		return false
	}
	this.cond.Signal()
	//ensure totally stopped
	this.cond.L.Lock()
	this.cond.L.Unlock()
	return true //success
}

//可以在不同的 goroutine 中执行 loop, 增加并行处理速度
func (this *EvtLoop) Loop() {
	for {
		this.cond.L.Lock()
		for {
			if this.IsClosed() {
				this.cond.L.Unlock()
				//退出执行
				return
			}
			if this.evtList.Len() != 0 {
				//不执行 cond 等待,直接取数据
				goto acquireDataAndProcess
			}
			if this.evtList.Len() == 0 {
				this.cond.Wait()
				continue
			}
		acquireDataAndProcess:
			evts := this.drainEvents()
			this.cond.L.Unlock()
			this.processEvent(evts)
			this.cond.L.Lock()
		}
	}
}

func (this *EvtLoop) drainEvents() *list.List {
	if this.evtList.Len() == 0 {
		return this.evtList
	}
	oldOne := this.evtList
	this.evtList = list.New()
	return oldOne
}

func (this *EvtLoop) processEvent(ies *list.List) {
	for it := ies.Front(); it != nil; it = it.Next() {
		//如果事件循环被打断,直接退出
		if this.IsClosed() {
			break
		}
		if it.Value != nil {
			this.evtHandler(it.Value)
		}
	}
}

func (this *EvtLoop) SendEventAsync(ie interface{}) {
	this.cond.L.Lock()
	needSignal := this.evtList.Len() == 0
	this.evtList.PushBack(ie)
	this.cond.L.Unlock()
	if needSignal {
		this.cond.Signal()
	}
}

func (this *EvtLoop) SendEventsAsync(ls *list.List) {
	this.cond.L.Lock()
	needSignal := this.evtList.Len() == 0
	this.evtList.PushBackList(ls)
	this.cond.L.Unlock()
	if needSignal {
		this.cond.Signal()
	}
}
