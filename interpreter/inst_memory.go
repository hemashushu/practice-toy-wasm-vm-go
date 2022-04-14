package interpreter

import (
	encoding_binary "encoding/binary"
	"errors"
	"wasmvm/binary"
)

// ======== 内存指令
//
// -------- 页面指令
//
// memory.size mem_block_idx:uint32
// 返回指定内存块的页面数（uint32）
//
// memory.grow mem_block_idx:uint32
// 增加指定数量的页面数
// 从操作数栈弹出 uint32 作为增加量
//
// 成功则返回旧的页面数量
// 失败（比如超出限制值的 max）则返回 -1
//
// 内存还有其他几个操作：
// - The `memory.fill` instruction sets all values in a region to a given byte.
// - The `memory.copy` instruction copies data from a source memory region to
//   a possibly overlapping destination region.
// - The `memory.init` instruction copies data from a passive data segment into a memory.
// - The `data.drop` instruction prevents further use of a passive data segment.
//   This instruction is intended to be used as an optimization hint.
//   After a data segment is dropped its data can no longer be retrieved,
//   so the memory used by this segment may be freed.
//
// https://webassembly.github.io/spec/core/syntax/instructions.html#syntax-instr-memory

func memorySize(v *vm, args interface{}) {
	mem_block_idx := args.(byte)
	if mem_block_idx != 0 {
		panic(errors.New("invalid memory index"))
	}

	v.operandStack.pushU32(v.memory.Size())
}

func memoryGrow(v *vm, args interface{}) {
	mem_block_idx := args.(byte)
	if mem_block_idx != 0 {
		panic(errors.New("invalid memory index"))
	}

	previousSize := v.memory.Grow(v.operandStack.popU32())
	// 虽然 grow 指令有可能会返回 -1，但仅表示失败，所以指令的返回值
	// 仍然以 u32 类型压入栈
	v.operandStack.pushU32(previousSize)
}

// ------ 加载指令
//
// i32.load align:uint32 offset:uint32
//
// load 指令有两个立即数：
// - align 地址对齐字节数量的对数，表示对齐一个 ”以 2 为底，以 align 为指数“ 的字节数，
//   比如 align = 1 时，表示对齐 2^1 = 2 个字节
//   比如 align = 2 时，表示对齐 2^2 = 4 个字节
//   align 只起提示作用，用于帮助编译器优化机器代码，对实际执行没有影响（对于 wasm 解析器，可以忽略这个值）
//
// - offset 偏移值
//   加载（以及存储）指令都会从操作数栈弹出一个 i32 类型的整数，让它与指令的立即数 offset 相加，得到
//   实际的内存地址，即：有效地址 = offset + popUint32()
//
// 加载过程：
// 1. 从操作数栈弹出一个 uint32，作为目标地址（addr）
// 2. 计算有效地址
// 3. 将指定地址内存的数值（多个字节使用小端格式 litte-endian 编码）压入操作数栈
//
// 指令列表
// i32.load
// i32.load16_s
// i32.load16_u
// i32.load8_s
// i32.load8_u
//
// i64.load
// i64.load32_s
// i64.load32_u
// i64.load16_s
// i64.load16_u
// i64.load8_s
// i64.load8_u
//
// f32.load
// f64.load

func i32Load(v *vm, memArg interface{}) {
	val := readU32(v, memArg)
	v.operandStack.pushU32(val)
}

func i32Load8S(v *vm, memArg interface{}) {
	val := readU8(v, memArg)
	v.operandStack.pushS32(int32(int8(val)))
}

func i32Load8U(v *vm, memArg interface{}) {
	val := readU8(v, memArg)
	v.operandStack.pushU32(uint32(val))
}

func i32Load16S(v *vm, memArg interface{}) {
	val := readU16(v, memArg)
	v.operandStack.pushS32(int32(int16(val)))
}

func i32Load16U(v *vm, memArg interface{}) {
	val := readU16(v, memArg)
	v.operandStack.pushU32(uint32(val))
}

func i64Load(v *vm, memArg interface{}) {
	val := readU64(v, memArg)
	v.operandStack.pushU64(val)
}

func i64Load8S(v *vm, memArg interface{}) {
	val := readU8(v, memArg)
	v.operandStack.pushS64(int64(int8(val)))
}

func i64Load8U(v *vm, memArg interface{}) {
	val := readU8(v, memArg)
	v.operandStack.pushU64(uint64(val))
}

func i64Load16S(v *vm, memArg interface{}) {
	val := readU16(v, memArg)
	v.operandStack.pushS64(int64(int16(val)))
}

func i64Load16U(v *vm, memArg interface{}) {
	val := readU16(v, memArg)
	v.operandStack.pushU64(uint64(val))
}

func i64Load32S(v *vm, memArg interface{}) {
	val := readU32(v, memArg)
	v.operandStack.pushS64(int64(int32(val)))
}

func i64Load32U(v *vm, memArg interface{}) {
	val := readU32(v, memArg)
	v.operandStack.pushU64(uint64(val))
}

func f32Load(v *vm, memArg interface{}) {
	val := readU32(v, memArg)
	v.operandStack.pushU32(val)
}

func f64Load(v *vm, memArg interface{}) {
	val := readU64(v, memArg)
	v.operandStack.pushU64(val)
}

var byteOrder = encoding_binary.LittleEndian

// 辅助函数

func readU8(v *vm, memArg interface{}) byte {
	var buf [1]byte
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Read(eaddr, buf[:])
	return buf[0]
}

func readU16(v *vm, memArg interface{}) uint16 {
	var buf [2]byte
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Read(eaddr, buf[:])
	return byteOrder.Uint16(buf[:])
}

func readU32(v *vm, memArg interface{}) uint32 {
	var buf [4]byte
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Read(eaddr, buf[:])
	return byteOrder.Uint32(buf[:])
}

func readU64(v *vm, memArg interface{}) uint64 {
	var buf [8]byte
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Read(eaddr, buf[:])
	return byteOrder.Uint64(buf[:])
}

// 因为指令中的 offset 立即数是 uint32，而操作数栈弹出的值也是 uint32，
// 所以有效地址（uint32 + uint32）是一个 33 位的无符号整数，超出了 uint32 的范围，
// 所以这里使用 uint64 存储有效地址。
func getEffectiveAddress(v *vm, memArg interface{}) uint64 {
	// MemArg 里头的 align 暂时无用
	realMemArg := memArg.(binary.MemArg)
	offset := realMemArg.Offset
	return uint64(v.operandStack.popU32()) + uint64(offset)
}

// -------- 存储指令
//
// i32.store align:uint32 offset:uint32
//
// 加载过程：
// 1. 从操作数栈弹出一个操作数，这个操作数将作为被存储的数据（data）
// 2. 从操作数栈弹出一个 uint32，作为目标地址（addr）
// 3. 计算有效地址
// 4. 将 data 写入指定地址的内存
//
// i32.store
// i32.store_16
// i32.store_8
//
// i64.store
// i64.store_32
// i64.store_16
// i64.store_8
//
// f32.store
// f64.store

func i32Store(v *vm, memArg interface{}) {
	val := v.operandStack.popU32()
	writeU32(v, memArg, val)
}

func i32Store8(v *vm, memArg interface{}) {
	val := v.operandStack.popU32()
	writeU8(v, memArg, byte(val))
}

func i32Store16(v *vm, memArg interface{}) {
	val := v.operandStack.popU32()
	writeU16(v, memArg, uint16(val))
}

func i64Store(v *vm, memArg interface{}) {
	val := v.operandStack.popU64()
	writeU64(v, memArg, val)
}

func i64Store8(v *vm, memArg interface{}) {
	val := v.operandStack.popU64()
	writeU8(v, memArg, byte(val))
}

func i64Store16(v *vm, memArg interface{}) {
	val := v.operandStack.popU64()
	writeU16(v, memArg, uint16(val))
}

func i64Store32(v *vm, memArg interface{}) {
	val := v.operandStack.popU64()
	writeU32(v, memArg, uint32(val))
}

func f32Store(v *vm, memArg interface{}) {
	val := v.operandStack.popU32()
	writeU32(v, memArg, val)
}

func f64Store(v *vm, memArg interface{}) {
	val := v.operandStack.popU64()
	writeU64(v, memArg, val)
}

// 辅助函数

func writeU8(v *vm, memArg interface{}, val byte) {
	var buf [1]byte
	buf[0] = val
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Write(eaddr, buf[:])
}

func writeU16(v *vm, memArg interface{}, n uint16) {
	var buf [2]byte
	byteOrder.PutUint16(buf[:], n)
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Write(eaddr, buf[:])
}

func writeU32(v *vm, memArg interface{}, n uint32) {
	var buf [4]byte
	byteOrder.PutUint32(buf[:], n)
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Write(eaddr, buf[:])
}

func writeU64(v *vm, memArg interface{}, n uint64) {
	var buf [8]byte
	byteOrder.PutUint64(buf[:], n)
	eaddr := getEffectiveAddress(v, memArg)
	v.memory.Write(eaddr, buf[:])
}
