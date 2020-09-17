package msync

import "time"
import "github.com/guestin/mob/merrors"

type Pool struct {
	c       chan interface{}
	sz      int
	eleGen  func() interface{}
	eleFree func(interface{})
}

func NewPool(len int,
	eleGen func() interface{},
	eleFree func(interface{})) *Pool {
	merrors.Assert(len > 0, "bad len")
	merrors.Assert(eleGen != nil, "eleGen cant be nil")

	pool := &Pool{
		c:       make(chan interface{}, len),
		sz:      len,
		eleGen:  eleGen,
		eleFree: eleFree,
	}
	for idx := 0; idx < len; idx++ {
		ele := eleGen()
		merrors.Assert(ele != nil, "generate nil ele")
		pool.c <- ele
	}
	return pool
}

func (this *Pool) Free() {
	for idx := 0; idx < this.sz; idx++ {
		ele := this.Take()
		if this.eleFree == nil {
			return
		}
		this.eleFree(ele)
	}
	close(this.c)
}

func (this *Pool) Take() interface{} {
	return <-this.c
}

func (this *Pool) TakeWithin(timeout time.Duration) interface{} {
	afterC := time.After(timeout)
	select {
	case <-afterC:
		return nil
	case ele := <-this.c:
		return ele
	}
}

func (this *Pool) Return(ele interface{}) {
	merrors.Assert(ele != nil, "return ele(nil) to pool")
	this.c <- ele
}
