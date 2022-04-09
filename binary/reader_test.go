package binary

import (
	"testing"
	"wasmvm/assert"
)

func TestReads(t *testing.T) {
	reader := wasmReader{data: []byte{
		0x01,
		0x02, 0x03, 0x04, 0x05,
		0x00, 0x00, 0xc0, 0x3f,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x3f,
		0xE5, 0x8E, 0x26, // https://en.wikipedia.org/wiki/LEB128#Unsigned_LEB128
		0xC0, 0xBB, 0x78, // https://en.wikipedia.org/wiki/LEB128#Signed_LEB128
		0xC0, 0xBB, 0x78,
		0x03, 0x01, 0x02, 0x03,
		0x03, 0x66, 0x6f, 0x6f,
	}}
	assert.AssertEqual(t, byte(0x01), reader.readByte())
	assert.AssertEqual(t, uint32(0x05040302), reader.readU32())
	assert.AssertEqual(t, float32(1.5), reader.readF32())
	assert.AssertEqual(t, 1.5, reader.readF64())
	assert.AssertEqual(t, uint32(624485), reader.readVarU32())
	assert.AssertEqual(t, int32(-123456), reader.readVarS32())
	assert.AssertEqual(t, int64(-123456), reader.readVarS64())
	assert.AssertSliceEqual(t, []byte{0x01, 0x02, 0x03}, reader.readBytes())
	assert.AssertEqual(t, "foo", reader.readName())
	assert.AssertEqual(t, 0, reader.remaining())
}
