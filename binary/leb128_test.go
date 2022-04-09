package binary

import (
	"testing"
	"wasmvm/assert"
)

func TestDecodeVarUint(t *testing.T) {
	data := []byte{
		0b1_011_1111,
		0b1_001_1111,
		0b1_000_1111,
		0b1_000_0111,
		0b1_000_0011,
		0b0_000_0001}
	testDecodeVarUint32(t, data[5:], 0b0000001, 1)
	testDecodeVarUint32(t, data[4:], 0b1_0000011, 2)
	testDecodeVarUint32(t, data[3:], 0b1_0000011_0000111, 3)
	testDecodeVarUint32(t, data[2:], 0b1_0000011_0000111_0001111, 4)
	testDecodeVarUint32(t, data[1:], 0b1_0000011_0000111_0001111_0011111, 5)
}

func TestDecodeVarInt(t *testing.T) {
	data := []byte{0xC0, 0xBB, 0x78}
	testDecodeVarInt32(t, data, int32(-123456), 3)
}

func testDecodeVarUint32(t *testing.T, data []byte, expected_value uint32, expected_bytes int) {
	value, bytes := decodeVarUint(data, 32)

	assert.AssertEqual(t, expected_value, uint32(value))
	assert.AssertEqual(t, expected_bytes, bytes)
}
func testDecodeVarInt32(t *testing.T, data []byte, expected_value int32, expected_bytes int) {
	value, bytes := decodeVarInt(data, 32)

	assert.AssertEqual(t, expected_value, int32(value))
	assert.AssertEqual(t, expected_bytes, bytes)
}
