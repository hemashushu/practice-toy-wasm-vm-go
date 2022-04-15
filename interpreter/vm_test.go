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
	slots0 := runFuncAndGetStack("test-vm-const.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{123}, slots0)

	slots1 := runFuncAndGetStack("test-vm-const.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{123, 456}, slots1)
}

func TestInstParametric(t *testing.T) {
	slots0 := runFuncAndGetStack("test-vm-parametric.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{100, 123}, slots0)

	slots1 := runFuncAndGetStack("test-vm-parametric.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{100, 456}, slots1)

	slots2 := runFuncAndGetStack("test-vm-parametric.wasm", 2)
	assert.AssertSliceEqual(t, []uint64{123}, slots2)

	slots3 := runFuncAndGetStack("test-vm-parametric.wasm", 3)
	assert.AssertTrue(t, len(slots3) == 0)

	slots4 := runFuncAndGetStack("test-vm-parametric.wasm", 4)
	assert.AssertSliceEqual(t, []uint64{100, 123}, slots4)
}

func TestInstEqz(t *testing.T) {
	slots0 := runFuncAndGetStack("test-vm-eqz.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{10, 1}, slots0)

	slots1 := runFuncAndGetStack("test-vm-eqz.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{10, 0}, slots1)

	slots2 := runFuncAndGetStack("test-vm-eqz.wasm", 2)
	assert.AssertSliceEqual(t, []uint64{10, 0}, slots2)
}

func TestInstCompare(t *testing.T) {
	assert.AssertSliceEqual(t, []uint64{10, 0}, runFuncAndGetStack("test-vm-compare.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 2))
	assert.AssertSliceEqual(t, []uint64{10, 0}, runFuncAndGetStack("test-vm-compare.wasm", 3))

	assert.AssertSliceEqual(t, []uint64{10, 0}, runFuncAndGetStack("test-vm-compare.wasm", 4))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 6))
	assert.AssertSliceEqual(t, []uint64{10, 0}, runFuncAndGetStack("test-vm-compare.wasm", 7))

	assert.AssertSliceEqual(t, []uint64{10, 0}, runFuncAndGetStack("test-vm-compare.wasm", 8))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 9))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 10))
	assert.AssertSliceEqual(t, []uint64{10, 0}, runFuncAndGetStack("test-vm-compare.wasm", 11))

	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 12))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 13))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 14))
	assert.AssertSliceEqual(t, []uint64{10, 1}, runFuncAndGetStack("test-vm-compare.wasm", 15))

	// 测试浮点数

	assert.AssertSliceEqual(t, []uint64{11, 0}, runFuncAndGetStack("test-vm-compare.wasm", 16))
	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 17))
	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 18))
	assert.AssertSliceEqual(t, []uint64{11, 0}, runFuncAndGetStack("test-vm-compare.wasm", 19))
	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 20))
	assert.AssertSliceEqual(t, []uint64{11, 0}, runFuncAndGetStack("test-vm-compare.wasm", 21))
	assert.AssertSliceEqual(t, []uint64{11, 0}, runFuncAndGetStack("test-vm-compare.wasm", 22))
	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 23))

	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 24))
	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 25))
	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 26))
	assert.AssertSliceEqual(t, []uint64{11, 1}, runFuncAndGetStack("test-vm-compare.wasm", 27))
}

func TestInstUnary(t *testing.T) {
	assert.AssertSliceEqual(t, []uint64{27}, runFuncAndGetStack("test-vm-unary.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{2}, runFuncAndGetStack("test-vm-unary.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{3}, runFuncAndGetStack("test-vm-unary.wasm", 2))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.718))}, runFuncAndGetStack("test-vm-unary.wasm", 3))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.718))}, runFuncAndGetStack("test-vm-unary.wasm", 4))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(-2.718))}, runFuncAndGetStack("test-vm-unary.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(3.0))}, runFuncAndGetStack("test-vm-unary.wasm", 6))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, runFuncAndGetStack("test-vm-unary.wasm", 7))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, runFuncAndGetStack("test-vm-unary.wasm", 8))

	// 就近取整（4 舍 6 入，5 奇进偶不进）
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(1.0))}, runFuncAndGetStack("test-vm-unary.wasm", 9))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, runFuncAndGetStack("test-vm-unary.wasm", 10))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(2.0))}, runFuncAndGetStack("test-vm-unary.wasm", 11))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(4.0))}, runFuncAndGetStack("test-vm-unary.wasm", 12))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(5.0))}, runFuncAndGetStack("test-vm-unary.wasm", 13))
}

func TestInstBinary(t *testing.T) {
	n1 := -11
	n2 := -4
	n3 := -2

	assert.AssertSliceEqual(t, []uint64{11, 55}, runFuncAndGetStack("test-vm-binary.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{11, uint64(uint32(n1))}, runFuncAndGetStack("test-vm-binary.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{11, 726}, runFuncAndGetStack("test-vm-binary.wasm", 2))
	assert.AssertSliceEqual(t, []uint64{11, uint64(uint32(n2))}, runFuncAndGetStack("test-vm-binary.wasm", 3))
	assert.AssertSliceEqual(t, []uint64{11, 0b01111111111111111111111111111100}, runFuncAndGetStack("test-vm-binary.wasm", 4))
	assert.AssertSliceEqual(t, []uint64{11, uint64(uint32(n3))}, runFuncAndGetStack("test-vm-binary.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{11, 2}, runFuncAndGetStack("test-vm-binary.wasm", 6))

	// 位运算
	assert.AssertSliceEqual(t, []uint64{11, 0b11000}, runFuncAndGetStack("test-vm-binary.wasm", 7))
	assert.AssertSliceEqual(t, []uint64{11, 0b1111_1001}, runFuncAndGetStack("test-vm-binary.wasm", 8))
	assert.AssertSliceEqual(t, []uint64{11, 0b1110_0001}, runFuncAndGetStack("test-vm-binary.wasm", 9))

	assert.AssertSliceEqual(t, []uint64{11, 0b11111111_11111111_11111111_1111_0000}, runFuncAndGetStack("test-vm-binary.wasm", 10))
	assert.AssertSliceEqual(t, []uint64{11, 0b11111111_11111111_11111111_1111_1111}, runFuncAndGetStack("test-vm-binary.wasm", 11))
	assert.AssertSliceEqual(t, []uint64{11, 0b00001111_11111111_11111111_1111_1111}, runFuncAndGetStack("test-vm-binary.wasm", 12))

	assert.AssertSliceEqual(t, []uint64{11, 0b11111111_11111111_11111111_1110_0011}, runFuncAndGetStack("test-vm-binary.wasm", 13))
	assert.AssertSliceEqual(t, []uint64{11, 0b00_11111111_11111111_11111111_111110}, runFuncAndGetStack("test-vm-binary.wasm", 14))

}

func TestInstConvert(t *testing.T) {
	n1 := -8

	assert.AssertSliceEqual(t, []uint64{123}, runFuncAndGetStack("test-vm-convert.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{8}, runFuncAndGetStack("test-vm-convert.wasm", 1))
	assert.AssertSliceEqual(t, []uint64{8}, runFuncAndGetStack("test-vm-convert.wasm", 2))
	assert.AssertSliceEqual(t, []uint64{uint64(n1)}, runFuncAndGetStack("test-vm-convert.wasm", 3))
	assert.AssertSliceEqual(t, []uint64{0x00_00_00_00_ff_ff_ff_f8}, runFuncAndGetStack("test-vm-convert.wasm", 4))

	assert.AssertSliceEqual(t, []uint64{3}, runFuncAndGetStack("test-vm-convert.wasm", 5))
	assert.AssertSliceEqual(t, []uint64{3}, runFuncAndGetStack("test-vm-convert.wasm", 6))

	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(66.0))}, runFuncAndGetStack("test-vm-convert.wasm", 7))
	assert.AssertSliceEqual(t, []uint64{uint64(math.Float32bits(66.0))}, runFuncAndGetStack("test-vm-convert.wasm", 8))

	// todo:: 这里仅测试了部分指令
}

func TestInstMemoryPage(t *testing.T) {
	assert.AssertSliceEqual(t, []uint64{10, 2}, runFuncAndGetStack("test-vm-memory-page.wasm", 0))
	assert.AssertSliceEqual(t, []uint64{10, 2, 4, 7}, runFuncAndGetStack("test-vm-memory-page.wasm", 1))
}

func TestInstMemoryLoad(t *testing.T) {
	var init_data []byte = []byte{
		/* addr: 0      */ 0x11, // 17
		/* addr: 1      */ 0xf1, // uint8'241 == int8'-15 (-15=241-256)
		/* addr: 2,3    */ 0x55, 0x66, // 0x6655
		/* addr: 4,5    */ 0x80, 0x90, // 0x9080
		/* addr: 6..13  */ 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		/* addr: 14..21 */ 0x80, 0x90, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0,
	}

	s0, _ := runFuncWithInitMemoryData("test-vm-memory-load.wasm", init_data, 0)
	assert.AssertSliceEqual(t, []uint64{0x11}, s0)

	s1, _ := runFuncWithInitMemoryData("test-vm-memory-load.wasm", init_data, 1)
	assert.AssertSliceEqual(t, []uint64{0x11, 0xf1, 0x55, 0x66}, s1)

	s2, _ := runFuncWithInitMemoryData("test-vm-memory-load.wasm", init_data, 2)
	assert.AssertSliceEqual(t, []uint64{0x11, 0xf1, 0x55, 0x66}, s2)

	// 测试符号
	s3, _ := runFuncWithInitMemoryData("test-vm-memory-load.wasm", init_data, 3)
	n1 := -15
	assert.AssertSliceEqual(t, []uint64{17, 17, 241, uint64(uint32(n1))}, s3)

	// 测试 16 位和 32 位整数
	s4, _ := runFuncWithInitMemoryData("test-vm-memory-load.wasm", init_data, 4)
	assert.AssertSliceEqual(t, []uint64{0x6655, 0x6655, 0x9080, 0xffff9080, 0x03020100}, s4)

	// 测试 64 位整数
	s5, _ := runFuncWithInitMemoryData("test-vm-memory-load.wasm", init_data, 5)
	assert.AssertSliceEqual(t, []uint64{
		0x03020100,
		0x03020100,
		0xb0a09080,
		0xffffffffb0a09080,
		0x0706050403020100,
		0xf0e0d0c0b0a09080}, s5)
}

func TestInstMemoryStore(t *testing.T) {
	// addr: 0      : 0x11,
	// addr: 8      : 0x2233				: 0x33, 0x22
	// addr: 16     : 0x44556677			: 0x77, 0x66, 0x55, 0x44
	// addr: 24     : 0xf0e0d0c0b0a09080	: 0x80, 0x90, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0
	// addr: 32     : "hello"  				: 0x68, 0x65, 0x6C, 0x6C, 0x6F,
	// addr: 40		: "中文"    			 : 0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87,

	s0, m0 := runFunc("test-vm-memory-store.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{
		0x11,
		0x2233,
		0x44556677,
		0xf0e0d0c0b0a09080,
		0x68,
		0xe4}, s0)

	assertPartialMemoryData(t, []byte{
		0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x33, 0x22, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x77, 0x66, 0x55, 0x44, 0x00, 0x00, 0x00, 0x00,
		0x80, 0x90, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0,
		0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00, 0x00, 0x00,
		0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87, 0x00, 0x00,
	}, m0)

	s1, m1 := runFunc("test-vm-memory-store.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{0xddccbbaa}, s1)
	assertPartialMemoryData(t, []byte{
		0xaa, 0xbb, 0xcc, 0xdd, 0x00, 0x00, 0x00, 0x00,
	}, m1)

	s2, m2 := runFunc("test-vm-memory-store.wasm", 2)
	assert.AssertSliceEqual(t, []uint64{}, s2)
	assertPartialMemoryData(t, []byte{
		0x11, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x00,
		0xa3, 0xa2, 0xa1, 0xa0, 0xb0, 0x00, 0xc1, 0xc0,
		0xd3, 0xd2, 0xd1, 0xd0,
		0xe7, 0xe6, 0xe5, 0xe4, 0xe3, 0xe2, 0xe1, 0xe0,
	}, m2)

}

func TestInstMemoryGrow(t *testing.T) {
	assert.AssertSliceEqual(t, []uint64{
		1,
		6,
		0xaabbccdd,
		0x10012002,
	}, runFuncAndGetStack("test-vm-memory-grow.wasm", 0))
}

func assertPartialMemoryData(t *testing.T, expected []byte, actual []byte) {
	partial := make([]byte, len(expected))
	copy(partial, actual)
	assert.AssertSliceEqual(t, expected, partial)
}

func runFuncAndGetStack(fileName string, func_idx uint32) []uint64 {
	stack, _ := runFunc(fileName, func_idx)
	return stack
}

func runFunc(fileName string, func_idx uint32) ([]uint64, []byte) {
	m := readModule(fileName)
	return TestFunc(m, func_idx)
}

func runFuncWithInitMemoryData(fileName string, init_data []byte, func_idx uint32) ([]uint64, []byte) {
	m := readModule(fileName)
	return TestFuncWithInitMemoryData(m, init_data, func_idx)
}

func readModule(fileName string) binary.Module {
	currentDir, err := os.Getwd() // Getwd() 返回当前 package 的目录，比如 `/path/to/project/binary`
	if err != nil {
		panic(err)
	}

	testResourcesDir := filepath.Join(currentDir, "..", "test", "resources", "interpreter")
	wasmFilePath := filepath.Join(testResourcesDir, fileName)

	return binary.DecodeFile(wasmFilePath)
}
