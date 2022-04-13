package interpreter

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"wasmvm/assert"
	"wasmvm/binary"
)

func TestInstConst(t *testing.T) {
	slots0 := run("test-vm-const.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{123}, slots0)

	slots1 := run("test-vm-const.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{123, 456}, slots1)
}

func TestInstOperand(t *testing.T) {
	slots0 := run("test-vm-operand.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{100, 123}, slots0)

	slots1 := run("test-vm-operand.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{100, 456}, slots1)

	slots2 := run("test-vm-operand.wasm", 2)
	assert.AssertSliceEqual(t, []uint64{123}, slots2)

	slots3 := run("test-vm-operand.wasm", 3)
	assert.AssertTrue(t, len(slots3) == 0)

	slots4 := run("test-vm-operand.wasm", 4)
	assert.AssertSliceEqual(t, []uint64{100, 123}, slots4)
}

func TestInstEqz(t *testing.T) {
	slots0 := run("test-vm-eqz.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{10, 1}, slots0)

	slots1 := run("test-vm-eqz.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{10, 0}, slots1)

	slots2 := run("test-vm-eqz.wasm", 2)
	assert.AssertSliceEqual(t, []uint64{10, 0}, slots2)
}

func TestInstCompare(t *testing.T) {
	assert.AssertSliceEqual(t, []uint64{10, 0}, run("test-vm-compare.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 2))
	assert.AssertSliceEqual(t, []uint64{10, 0}, run("test-vm-compare.wasm", 3))

	assert.AssertSliceEqual(t, []uint64{10, 0}, run("test-vm-compare.wasm", 4))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 6))
	assert.AssertSliceEqual(t, []uint64{10, 0}, run("test-vm-compare.wasm", 7))

	assert.AssertSliceEqual(t, []uint64{10, 0}, run("test-vm-compare.wasm", 8))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 9))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 10))
	assert.AssertSliceEqual(t, []uint64{10, 0}, run("test-vm-compare.wasm", 11))

	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 12))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 13))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 14))
	assert.AssertSliceEqual(t, []uint64{10, 1}, run("test-vm-compare.wasm", 15))

	// 测试浮点数

	assert.AssertSliceEqual(t, []uint64{11, 0}, run("test-vm-compare.wasm", 16))
	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 17))
	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 18))
	assert.AssertSliceEqual(t, []uint64{11, 0}, run("test-vm-compare.wasm", 19))
	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 20))
	assert.AssertSliceEqual(t, []uint64{11, 0}, run("test-vm-compare.wasm", 21))
	assert.AssertSliceEqual(t, []uint64{11, 0}, run("test-vm-compare.wasm", 22))
	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 23))

	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 24))
	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 25))
	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 26))
	assert.AssertSliceEqual(t, []uint64{11, 1}, run("test-vm-compare.wasm", 27))

}

func TestInstUnary(t *testing.T) {
	assert.AssertSliceEqual(t, []uint64{27}, run("test-vm-unary.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{2}, run("test-vm-unary.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{3}, run("test-vm-unary.wasm", 2))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.718))}, run("test-vm-unary.wasm", 3))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.718))}, run("test-vm-unary.wasm", 4))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(-2.718))}, run("test-vm-unary.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(3.0))}, run("test-vm-unary.wasm", 6))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, run("test-vm-unary.wasm", 7))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, run("test-vm-unary.wasm", 8))

	// 就近取整（4 舍 6 入，5 奇进偶不进）
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(1.0))}, run("test-vm-unary.wasm", 9))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, run("test-vm-unary.wasm", 10))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, run("test-vm-unary.wasm", 11))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(4.0))}, run("test-vm-unary.wasm", 12))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(5.0))}, run("test-vm-unary.wasm", 13))
}

func TestInstBinary(t *testing.T) {
	n1 := -11
	n2 := -4
	n3 := -2

	assert.AssertSliceEqual(t, []uint64{11, 55}, run("test-vm-binary.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{11, uint64(uint32(n1))}, run("test-vm-binary.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{11, 726}, run("test-vm-binary.wasm", 2))
	assert.AssertSliceEqual(t, []uint64{11, uint64(uint32(n2))}, run("test-vm-binary.wasm", 3))
	assert.AssertSliceEqual(t, []uint64{11, 0b01111111111111111111111111111100}, run("test-vm-binary.wasm", 4))
	assert.AssertSliceEqual(t, []uint64{11, uint64(uint32(n3))}, run("test-vm-binary.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{11, 2}, run("test-vm-binary.wasm", 6))

	// 位运算
	assert.AssertSliceEqual(t, []uint64{11, 0b11000}, run("test-vm-binary.wasm", 7))
	assert.AssertSliceEqual(t, []uint64{11, 0b1111_1001}, run("test-vm-binary.wasm", 8))
	assert.AssertSliceEqual(t, []uint64{11, 0b1110_0001}, run("test-vm-binary.wasm", 9))

	assert.AssertSliceEqual(t, []uint64{11, 0b11111111_11111111_11111111_1111_0000}, run("test-vm-binary.wasm", 10))
	assert.AssertSliceEqual(t, []uint64{11, 0b11111111_11111111_11111111_1111_1111}, run("test-vm-binary.wasm", 11))
	assert.AssertSliceEqual(t, []uint64{11, 0b00001111_11111111_11111111_1111_1111}, run("test-vm-binary.wasm", 12))

	assert.AssertSliceEqual(t, []uint64{11, 0b11111111_11111111_11111111_1110_0011}, run("test-vm-binary.wasm", 13))
	assert.AssertSliceEqual(t, []uint64{11, 0b00_11111111_11111111_11111111_111110}, run("test-vm-binary.wasm", 14))

}

func TestInstConvert(t *testing.T) {
	n1 := -8

	assert.AssertSliceEqual(t, []uint64{123}, run("test-vm-convert.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{8}, run("test-vm-convert.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{8}, run("test-vm-convert.wasm", 2))
	assert.AssertSliceEqual(t, []uint64{uint64(n1)}, run("test-vm-convert.wasm", 3))
	assert.AssertSliceEqual(t, []uint64{0x00_00_00_00_ff_ff_ff_f8}, run("test-vm-convert.wasm", 4))

	assert.AssertSliceEqual(t, []uint64{3}, run("test-vm-convert.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{3}, run("test-vm-convert.wasm", 6))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(66.0))}, run("test-vm-convert.wasm", 7))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(66.0))}, run("test-vm-convert.wasm", 8))

}

func run(fileName string, func_idx uint32) []uint64 {
	currentDir, err := os.Getwd() // Getwd() 返回当前 package 的目录，比如 `/path/to/project/binary`
	if err != nil {
		panic(err)
	}

	testResourcesDir := filepath.Join(currentDir, "..", "test", "resources", "interpreter")
	wasmFilePath := filepath.Join(testResourcesDir, fileName)

	m := binary.DecodeFile(wasmFilePath)
	return ExecFunc(m, func_idx)
}
