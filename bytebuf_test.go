package mod

import (
	"github.com/guestin/mob/merrors"
	"github.com/guestin/mob/mio"
	"testing"
)

func TestByteBuffer(t *testing.T) {
	buffer := mio.NewByteBuffer()
	buffer.Write([]byte{0x01, 0x0e})
	merrors.Assert(buffer.PeekInt8(1) == int8(0x0e), "bad")
	merrors.Assert(buffer.ReadInt16() == 0x010e, "bad")
	merrors.Assert(buffer.ReadableLen() == 0, "bad")

	buffer.WriteUInt16(0x010e)
	merrors.Assert(buffer.PeekInt8(1) == int8(0x0e), "bad")
	merrors.Assert(buffer.ReadInt16() == 0x010e, "bad")
	merrors.Assert(buffer.ReadableLen() == 0, "bad")

	buffer.WriteInt32(0x010e)
	merrors.Assert(buffer.PeekInt8(buffer.ReadableLen()-1) == int8(0x0e), "bad")
	merrors.Assert(buffer.ReadInt32() == 0x010e, "bad")
	merrors.Assert(buffer.ReadableLen() == 0, "bad")
}
