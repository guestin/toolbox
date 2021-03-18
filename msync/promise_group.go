package msync

import (
	"github.com/guestin/mob/merrors"
	"sync"
)

type PromiseId string

type PromiseGroup struct {
	hangs map[PromiseId]Promise
	lock  sync.Locker
}

func NewPromiseGroup() *PromiseGroup {
	ret := &PromiseGroup{
		hangs: make(map[PromiseId]Promise),
		lock:  &sync.Mutex{},
	}
	return ret
}

func (this *PromiseGroup) DonePromise(id PromiseId, err error, data interface{}) {
	defer this.lock.Unlock()
	this.lock.Lock()
	if c, ok := this.hangs[id]; ok {
		c.DoneData(err, data)
	}
	delete(this.hangs, id)
}

var ErrPromiseReplace = merrors.ErrorWrap(nil, "this promise has been replace by another one!")

func (this *PromiseGroup) AddPromise(id PromiseId, c Promise) {
	defer this.lock.Unlock()
	this.lock.Lock()
	old, ok := this.hangs[id]
	if !ok {
		this.hangs[id] = c
		return
	}
	old.Done(ErrPromiseReplace)
}

func (this *PromiseGroup) RemovePromise(id PromiseId) {
	defer this.lock.Unlock()
	this.lock.Lock()
	delete(this.hangs, id)
}
