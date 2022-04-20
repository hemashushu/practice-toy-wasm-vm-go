package interpreter

import (
	"os"
	"path/filepath"
	"testing"
	"wasmvm/assert"
	"wasmvm/binary"
)

func wrapList[T comparable](items []T) []interface{} {
	rs := []interface{}{}
	for _, i := range items {
		rs = append(rs, i)
	}
	return rs
}

func TestInstConst(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{123}), testFuncWithoutArgs("test-vm-const.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{123, 456}), testFuncWithoutArgs("test-vm-const.wasm", 1))
}

func TestInstParametric(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{100, 123}), testFuncWithoutArgs("test-vm-parametric.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{100, 456}), testFuncWithoutArgs("test-vm-parametric.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{123}), testFuncWithoutArgs("test-vm-parametric.wasm", 2))
	assert.AssertListEqual(t, wrapList([]int32{}), testFuncWithoutArgs("test-vm-parametric.wasm", 3))
	assert.AssertListEqual(t, wrapList([]int32{100, 123}), testFuncWithoutArgs("test-vm-parametric.wasm", 4))
}

func TestInstEqz(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-eqz.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-eqz.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-eqz.wasm", 2))
}

func TestInstCompare(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 2))
	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 3))

	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 4))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 5))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 6))
	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 7))

	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 8))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 9))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 10))
	assert.AssertListEqual(t, wrapList([]int32{10, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 11))

	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 12))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 13))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 14))
	assert.AssertListEqual(t, wrapList([]int32{10, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 15))

	// 测试浮点数的比较

	assert.AssertListEqual(t, wrapList([]int32{11, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 16))
	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 17))
	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 18))
	assert.AssertListEqual(t, wrapList([]int32{11, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 19))
	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 20))
	assert.AssertListEqual(t, wrapList([]int32{11, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 21))
	assert.AssertListEqual(t, wrapList([]int32{11, 0}), testFuncWithoutArgs("test-vm-compare.wasm", 22))
	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 23))

	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 24))
	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 25))
	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 26))
	assert.AssertListEqual(t, wrapList([]int32{11, 1}), testFuncWithoutArgs("test-vm-compare.wasm", 27))
}

func TestInstUnary(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{27}), testFuncWithoutArgs("test-vm-unary.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{2}), testFuncWithoutArgs("test-vm-unary.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{3}), testFuncWithoutArgs("test-vm-unary.wasm", 2))

	assert.AssertListEqual(t, wrapList([]float32{2.718}), testFuncWithoutArgs("test-vm-unary.wasm", 3))
	assert.AssertListEqual(t, wrapList([]float32{2.718}), testFuncWithoutArgs("test-vm-unary.wasm", 4))
	assert.AssertListEqual(t, wrapList([]float32{-2.718}), testFuncWithoutArgs("test-vm-unary.wasm", 5))
	assert.AssertListEqual(t, wrapList([]float32{3.0}), testFuncWithoutArgs("test-vm-unary.wasm", 6))
	assert.AssertListEqual(t, wrapList([]float32{2.0}), testFuncWithoutArgs("test-vm-unary.wasm", 7))

	assert.AssertListEqual(t, wrapList([]float32{2.0}), testFuncWithoutArgs("test-vm-unary.wasm", 8))

	// 就近取整（4 舍 6 入，5 奇进偶不进）
	assert.AssertListEqual(t, wrapList([]float32{1.0}), testFuncWithoutArgs("test-vm-unary.wasm", 9))
	assert.AssertListEqual(t, wrapList([]float32{2.0}), testFuncWithoutArgs("test-vm-unary.wasm", 10))
	assert.AssertListEqual(t, wrapList([]float32{2.0}), testFuncWithoutArgs("test-vm-unary.wasm", 11))
	assert.AssertListEqual(t, wrapList([]float32{4.0}), testFuncWithoutArgs("test-vm-unary.wasm", 12))

	assert.AssertListEqual(t, wrapList([]float32{5.0}), testFuncWithoutArgs("test-vm-unary.wasm", 13))
}

func TestInstBinary(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{11, 55}), testFuncWithoutArgs("test-vm-binary.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{11, -11 /*uint64(uint32(n1))*/}), testFuncWithoutArgs("test-vm-binary.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{11, 726}), testFuncWithoutArgs("test-vm-binary.wasm", 2))
	assert.AssertListEqual(t, wrapList([]int32{11, -4 /*uint64(uint32(n2))*/}), testFuncWithoutArgs("test-vm-binary.wasm", 3))
	assert.AssertListEqual(t, wrapList([]int32{11, 0b01111111111111111111111111111100}), testFuncWithoutArgs("test-vm-binary.wasm", 4))
	assert.AssertListEqual(t, wrapList([]int32{11, -2 /*uint64(uint32(n3))*/}), testFuncWithoutArgs("test-vm-binary.wasm", 5))
	assert.AssertListEqual(t, wrapList([]int32{11, 2}), testFuncWithoutArgs("test-vm-binary.wasm", 6))

	// 位运算
	assert.AssertListEqual(t, wrapList([]int32{11, 0b11000}), testFuncWithoutArgs("test-vm-binary.wasm", 7))
	assert.AssertListEqual(t, wrapList([]int32{11, 0b1111_1001}), testFuncWithoutArgs("test-vm-binary.wasm", 8))
	assert.AssertListEqual(t, wrapList([]int32{11, 0b1110_0001}), testFuncWithoutArgs("test-vm-binary.wasm", 9))

	n1 := 0b11111111_11111111_11111111_1111_0000
	assert.AssertListEqual(t, wrapList([]int32{11, int32(n1)}), testFuncWithoutArgs("test-vm-binary.wasm", 10))

	n2 := 0b11111111_11111111_11111111_1111_1111
	assert.AssertListEqual(t, wrapList([]int32{11, int32(n2)}), testFuncWithoutArgs("test-vm-binary.wasm", 11))
	assert.AssertListEqual(t, wrapList([]int32{11, 0b00001111_11111111_11111111_1111_1111}), testFuncWithoutArgs("test-vm-binary.wasm", 12))

	n3 := 0b11111111_11111111_11111111_1110_0011
	assert.AssertListEqual(t, wrapList([]int32{11, int32(n3)}), testFuncWithoutArgs("test-vm-binary.wasm", 13))
	assert.AssertListEqual(t, wrapList([]int32{11, 0b00_11111111_11111111_11111111_111110}), testFuncWithoutArgs("test-vm-binary.wasm", 14))

}

func TestInteroperate(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{22, 11}), testFunc("test-vm-interoperate.wasm", 0, wrapList([]int32{11, 22})))
	assert.AssertListEqual(t, wrapList([]int32{33}), testFunc("test-vm-interoperate.wasm", 1, wrapList([]int32{11, 22})))
}

func TestInstConvert(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{123}), testFuncWithoutArgs("test-vm-convert.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int64{8}), testFuncWithoutArgs("test-vm-convert.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int64{8}), testFuncWithoutArgs("test-vm-convert.wasm", 2))
	assert.AssertListEqual(t, wrapList([]int64{-8}), testFuncWithoutArgs("test-vm-convert.wasm", 3))
	assert.AssertListEqual(t, wrapList([]int64{0x00_00_00_00_ff_ff_ff_f8}), testFuncWithoutArgs("test-vm-convert.wasm", 4))

	assert.AssertListEqual(t, wrapList([]int32{3}), testFuncWithoutArgs("test-vm-convert.wasm", 5))
	assert.AssertListEqual(t, wrapList([]int32{3}), testFuncWithoutArgs("test-vm-convert.wasm", 6))

	assert.AssertListEqual(t, wrapList([]float32{66.0}), testFuncWithoutArgs("test-vm-convert.wasm", 7))
	assert.AssertListEqual(t, wrapList([]float32{66.0}), testFuncWithoutArgs("test-vm-convert.wasm", 8))

	// todo:: 这里仅测试了部分指令
}

func TestInstMemoryPage(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{10, 2}), testFuncWithoutArgs("test-vm-memory-page.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{10, 2, 4, 7}), testFuncWithoutArgs("test-vm-memory-page.wasm", 1))
}

func TestInstMemoryLoad(t *testing.T) {
	var init_memory_data []byte = []byte{
		/* addr: 0      */ 0x11, // 17
		/* addr: 1      */ 0xf1, // uint8'241 == int8'-15 (-15=241-256)
		/* addr: 2,3    */ 0x55, 0x66, // 0x6655
		/* addr: 4,5    */ 0x80, 0x90, // 0x9080
		/* addr: 6..13  */ 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		/* addr: 14..21 */ 0x80, 0x90, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0,
	}

	assert.AssertListEqual(t, wrapList([]int32{0x11}), testFuncWithInitMemoryData("test-vm-memory-load.wasm", init_memory_data, 0, nil))
	assert.AssertListEqual(t, wrapList([]int32{0x11, 0xf1, 0x55, 0x66}), testFuncWithInitMemoryData("test-vm-memory-load.wasm", init_memory_data, 1, nil))
	assert.AssertListEqual(t, wrapList([]int32{0x11, 0xf1, 0x55, 0x66}), testFuncWithInitMemoryData("test-vm-memory-load.wasm", init_memory_data, 2, nil))

	// 测试符号
	assert.AssertListEqual(t, wrapList([]int32{17, 17, 241, -15}), testFuncWithInitMemoryData("test-vm-memory-load.wasm", init_memory_data, 3, nil))

	// 测试 16 位和 32 位整数
	n1 := 0xffff9080
	assert.AssertListEqual(t, wrapList([]int32{0x6655, 0x6655, 0x9080, int32(n1), 0x03020100}), testFuncWithInitMemoryData("test-vm-memory-load.wasm", init_memory_data, 4, nil))

	// 测试 64 位整数
	var n2 uint64 = 0xffffffffb0a09080
	var n3 uint64 = 0xf0e0d0c0b0a09080
	assert.AssertListEqual(t, wrapList([]int64{
		0x03020100,
		0x03020100,
		0xb0a09080,
		int64(n2),
		0x0706050403020100,
		int64(n3)}), testFuncWithInitMemoryData("test-vm-memory-load.wasm", init_memory_data, 5, nil))
}

func TestInstMemoryStore(t *testing.T) {
	// addr: 0      : 0x11,
	// addr: 8      : 0x2233				: 0x33, 0x22
	// addr: 16     : 0x44556677			: 0x77, 0x66, 0x55, 0x44
	// addr: 24     : 0xf0e0d0c0b0a09080	: 0x80, 0x90, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0
	// addr: 32     : "hello"  				: 0x68, 0x65, 0x6C, 0x6C, 0x6F,
	// addr: 40		: "中文"    			 : 0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87,

	r0, m0 := testFuncAndDumpMemory("test-vm-memory-store.wasm", 0, nil)
	assert.AssertListEqual(t, wrapList([]int32{
		0x11,
		0x2233,
		0x44556677,
	}), r0)
	assertPartialMemoryData(t, []byte{
		0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x33, 0x22, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x77, 0x66, 0x55, 0x44, 0x00, 0x00, 0x00, 0x00,
		0x80, 0x90, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0,
		0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00, 0x00, 0x00,
		0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87, 0x00, 0x00,
	}, m0)

	r1 := testFunc("test-vm-memory-store.wasm", 1, nil)
	var n1 uint64 = 0xf0e0d0c0b0a09080
	assert.AssertListEqual(t, wrapList([]int64{
		int64(n1),
		0x68,
		0xe4}), r1)

	r2, m2 := testFuncAndDumpMemory("test-vm-memory-store.wasm", 2, nil)
	var n2 uint32 = 0xddccbbaa
	assert.AssertListEqual(t,
		wrapList([]int32{int32(n2)}), r2)
	assertPartialMemoryData(t, []byte{
		0xaa, 0xbb, 0xcc, 0xdd, 0x00, 0x00, 0x00, 0x00,
	}, m2)

	r3, m3 := testFuncAndDumpMemory("test-vm-memory-store.wasm", 3, nil)
	assert.AssertListEqual(t, []interface{}{}, r3)
	assertPartialMemoryData(t, []byte{
		0x11, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x00,
		0xa3, 0xa2, 0xa1, 0xa0, 0xb0, 0x00, 0xc1, 0xc0,
		0xd3, 0xd2, 0xd1, 0xd0,
		0xe7, 0xe6, 0xe5, 0xe4, 0xe3, 0xe2, 0xe1, 0xe0,
	}, m3)
}

func TestInstMemoryGrow(t *testing.T) {
	var n1 uint32 = 0xaabbccdd
	assert.AssertListEqual(t, wrapList([]int32{
		1,
		6,
		int32(n1),
		0x10012002,
	}), testFuncWithoutArgs("test-vm-memory-grow.wasm", 0))
}

func TestInstFunctionCall(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{3}), testFuncWithoutArgs("test-vm-function-call.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{1}), testFuncWithoutArgs("test-vm-function-call.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{-5}), testFuncWithoutArgs("test-vm-function-call.wasm", 2))
}

func TestInstControl(t *testing.T) {
	// 测试 return
	assert.AssertListEqual(t, wrapList([]int32{1}), testFuncWithoutArgs("test-vm-control.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{2}), testFuncWithoutArgs("test-vm-control.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{3}), testFuncWithoutArgs("test-vm-control.wasm", 2))

	// 测试 br
	assert.AssertListEqual(t, wrapList([]int32{4}), testFuncWithoutArgs("test-vm-control.wasm", 3))
	assert.AssertListEqual(t, wrapList([]int32{2}), testFuncWithoutArgs("test-vm-control.wasm", 4))
	assert.AssertListEqual(t, wrapList([]int32{11}), testFuncWithoutArgs("test-vm-control.wasm", 5))
	assert.AssertListEqual(t, wrapList([]int32{12}), testFuncWithoutArgs("test-vm-control.wasm", 6))
	assert.AssertListEqual(t, wrapList([]int32{13}), testFuncWithoutArgs("test-vm-control.wasm", 7))

	// 测试 br_if
	assert.AssertListEqual(t, wrapList([]int32{55}), testFuncWithoutArgs("test-vm-control.wasm", 8))

	// 测试 br_table
	// todo::

	// 测试 if
	assert.AssertListEqual(t, wrapList([]int32{2}), testFuncWithoutArgs("test-vm-control.wasm", 9))
	assert.AssertListEqual(t, wrapList([]int32{1}), testFuncWithoutArgs("test-vm-control.wasm", 10))
}

// func TestNativeFunction(t *testing.T) {
// 	// 测试调用本地函数（native function）
//
// 	// // 导入了 3 个 native function，所以内部函数的
// 	// // 索引值从 3 开始。
// 	// runFunc("test-vm-native-function.wasm", 3) // 用人类的眼睛观察输出窗口是否输出 "A"
// 	// runFunc("test-vm-native-function.wasm", 4) // 用人类的眼睛观察输出窗口是否输出 "65"
//
// 	assert.AssertListEqual(t, wrapList([]int32{33}, testFuncWithoutArgs("test-vm-native-function.wasm", 1))
// }

func TestInstCallIndirect(t *testing.T) {
	assert.AssertListEqual(t, wrapList([]int32{12}), testFuncWithoutArgs("test-vm-indirect-call.wasm", 0))
	assert.AssertListEqual(t, wrapList([]int32{8}), testFuncWithoutArgs("test-vm-indirect-call.wasm", 1))
	assert.AssertListEqual(t, wrapList([]int32{20}), testFuncWithoutArgs("test-vm-indirect-call.wasm", 2))
	assert.AssertListEqual(t, wrapList([]int32{5}), testFuncWithoutArgs("test-vm-indirect-call.wasm", 3))
}

// 仅测试实际数据当中部分的数据
func assertPartialMemoryData(t *testing.T, expected []byte, actual []byte) {
	partial := make([]byte, len(expected))
	copy(partial, actual)
	assert.AssertSliceEqual(t, expected, partial)
}

func testFuncWithoutArgs(fileName string, func_idx uint32) []interface{} {
	return testFunc(fileName, func_idx, nil)
}

func testFunc(fileName string, func_idx uint32, args []interface{}) []interface{} {
	m := readModule(fileName)
	return evalModuleFunc(m, func_idx, args)
}

func testFuncWithInitMemoryData(fileName string, init_memory_data []byte, func_idx uint32, args []interface{}) []interface{} {
	m := readModule(fileName)
	return evalModuleFuncWithInitMemoryData(m, init_memory_data, func_idx, args)
}

func testFuncAndDumpMemory(fileName string, func_idx uint32, args []interface{}) ([]interface{}, []byte) {
	m := readModule(fileName)
	return evalModuleFuncAndDumpMemory(m, func_idx, args)
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
