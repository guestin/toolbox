package mio

import (
	"container/list"
)

type RwBuffer interface {
	ReadableBuffer
	WritableBuffer
}

type ReadableBuffer interface {
	ReadInt32() int32
	ReadUInt32() uint32
	ReadInt16() int16
	ReadUInt16() uint16
	ReadInt8() int8
	ReadUInt8() uint8

	PeekInt32(offset int) int32
	PeekUInt32(offset int) uint32
	PeekInt16(offset int) int16
	PeekUInt16(offset int) uint16
	PeekInt8(offset int) int8
	PeekUInt8(offset int) uint8

	ReadableLen() int

	ReadN(n int) []uint8
	PeekN(offset, n int) []uint8
	PopN(n int)
	ToArray() []uint8
	PopAll()
}

type WritableBuffer interface {
	Write(data []byte)
	WriteUInt8(v uint8)
	WriteUInt16(v uint16)
	WriteInt32(v int32)
	WriteUInt32(v uint32)
}

type byteBufferImpl struct {
	data *list.List
}

func NewByteBuffer() RwBuffer {
	return &byteBufferImpl{
		data: list.New(),
	}
}

func (this *byteBufferImpl) ToArray() []uint8 {
	ret := make([]uint8, 0, this.ReadableLen())
	for ele := this.data.Front(); ele != nil; ele = ele.Next() {
		ret = append(ret, ele.Value.(uint8))
	}
	return ret
}

func (this *byteBufferImpl) PopAll() {
	this.PopN(this.ReadableLen())
}

func (this *byteBufferImpl) Write(arr []uint8) {
	for _, v := range arr {
		this.data.PushBack(v)
	}
}

func (this *byteBufferImpl) WriteUInt8(v uint8) {
	this.data.PushBack(v)
}

//checkout ReadableLen before call this
func (this *byteBufferImpl) PeekN(offset, n int) []uint8 {
	arr := make([]uint8, 0, n)
	var ele = this.data.Front()
	for i := 0; i < offset; ele = ele.Next() {
		if ele == nil {
			return arr
		}
		i++
	}
	for i, it := 0, ele; i < n && it != nil; it = it.Next() {
		arr = append(arr, it.Value.(uint8))
		i++
	}
	return arr
}

//checkout ReadableLen before call this
func (this *byteBufferImpl) ReadN(n int) []uint8 {
	arr := make([]uint8, 0, n)
	dataN := this.data.Len()
	for i := 0; i < n && i < dataN; i++ {
		front := this.data.Front()
		if front == nil {
			break
		}
		arr = append(arr, front.Value.(uint8))
		this.data.Remove(front)
	}
	return arr
}

//array len
func (this *byteBufferImpl) ReadableLen() int {
	return this.data.Len()
}

func (this *byteBufferImpl) PopN(n int) {
	for i := 0; i < n; i++ {
		front := this.data.Front()
		if front != nil {
			this.data.Remove(this.data.Front())
		}
	}
}

func (this *byteBufferImpl) ReadInt32() int32 {
	return ArrayToInt32BE(this.ReadN(4))
}

func (this *byteBufferImpl) ReadUInt32() uint32 {
	return ArrToUint32BE(this.ReadN(4))
}

func (this *byteBufferImpl) ReadInt16() int16 {
	return ArrToInt16BE(this.ReadN(2))
}

func (this *byteBufferImpl) ReadUInt16() uint16 {
	return ArrToUint16BE(this.ReadN(2))
}

func (this *byteBufferImpl) ReadInt8() int8 {
	return int8(this.ReadN(1)[0])
}

func (this *byteBufferImpl) ReadUInt8() uint8 {
	return this.ReadN(1)[0]
}

func (this *byteBufferImpl) PeekInt32(offset int) int32 {
	return ArrayToInt32BE(this.PeekN(offset, 4))
}

func (this *byteBufferImpl) PeekUInt32(offset int) uint32 {
	return ArrToUint32BE(this.PeekN(offset, 4))
}

func (this *byteBufferImpl) PeekInt16(offset int) int16 {
	return ArrToInt16BE(this.PeekN(offset, 2))
}

func (this *byteBufferImpl) PeekUInt16(offset int) uint16 {
	return ArrToUint16BE(this.PeekN(offset, 2))
}

func (this *byteBufferImpl) PeekInt8(offset int) int8 {
	return int8(this.PeekN(offset, 1)[0])
}

func (this *byteBufferImpl) PeekUInt8(offset int) uint8 {
	return this.PeekN(offset, 1)[0]
}

func (this *byteBufferImpl) WriteUInt16(v uint16) {
	this.Write(Uint16ToArrBE(v))
}

func (this *byteBufferImpl) WriteInt32(v int32) {
	this.Write(Int32ToArrBE(v))
}

func (this *byteBufferImpl) WriteUInt32(v uint32) {
	this.Write(UInt32ToArrBE(v))
}
