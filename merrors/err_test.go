package merrors

import (
	"errors"
	"reflect"
	"testing"
	"unsafe"
)

func TestErrorf(t *testing.T) {
	t.Log(Errorf("sample %s %s", "1", "2"))
}

func TestErrorWrapf(t *testing.T) {
	t.Log(ErrorWrapf(errors.New("raw error"), "sample %d/%d", 1, 2).SetData(1))
}

func TestErrorWrap(t *testing.T) {
	t.Log(ErrorWrap(errors.New("raw error"), "what fuck?"))
}

func TestEndianCast(t *testing.T) {
	//t.Log(Error("what fuck?"))
	v := uint16(0xFF00)
	slice := reflect.New(reflect.TypeOf(([]uint8)(nil)).Elem())
	header := (*reflect.SliceHeader)(unsafe.Pointer(slice.Pointer()))
	header.Cap = 2
	header.Len = 2
	header.Data = uintptr(unsafe.Pointer(&v))
	xxx := *(*[]uint8)(unsafe.Pointer(header))
	for _, v := range xxx {
		t.Logf("%02X ", v)
	}
	t.Log("------")
	v1 := uint8(v >> 8)
	v2 := uint8(v)
	t.Logf("%02X %02X", v1, v2)
}
