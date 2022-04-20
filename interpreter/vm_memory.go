package interpreter

import (
	"errors"
	"wasmvm/binary"
)

type memory struct {
	type_ binary.MemType // 限制值
	data  []byte         // 内存就是一个 byte 数组
}

func newMemory(memType binary.MemType) *memory {
	return &memory{
		type_: memType,
		data:  make([]byte, memType.Min*binary.PageSize),
	}
}

// 使用初始数据来创建内存块，用于测试
func newMemoryWithInitData(init_data []byte) *memory {
	data := make([]byte, 1*binary.PageSize) // 创建只有一个页面的内存块
	copy(data, init_data)
	return &memory{
		type_: binary.MemType{Tag: 0, Min: 1},
		data:  data,
	}
}

// 这是一个 Memory 接口需要的方法
func (m *memory) Type() binary.MemType {
	return m.type_
}

// 获取内存块的页面数
// 返回当前的页面数（uint32）
func (m *memory) Size() uint32 {
	return uint32(len(m.data) / binary.PageSize)
}

// 扩充内存大小（在内存块的 max 允许的范围之内）
// 参数 increaseCount: 需要增加的页面数，而不是 `增加到` 的页面数
// 返回旧的页面数（uint32）
// 失败时会返回被转为 uint32 的 -1
func (m *memory) Grow(increaseCount uint32) uint32 {
	previousSize := m.Size()
	if increaseCount == 0 {
		return previousSize
	}

	// 如果不指定 max 值，则可以增长到最大允许的页面数 MaxPageCount
	maxPages := uint32(binary.MaxPageCount)

	if m.type_.Max > 0 {
		maxPages = m.type_.Max
	}

	if previousSize+increaseCount > maxPages {
		// 失败则返回 -1
		n1 := -1
		return uint32(n1)
	}

	newData := make([]byte, (previousSize+increaseCount)*binary.PageSize)
	copy(newData, m.data)
	m.data = newData

	// 成功则返回之前的页面数
	// https://webassembly.github.io/spec/core/syntax/instructions.html#syntax-instr-memory
	// https://developer.mozilla.org/en-US/docs/WebAssembly/Reference/Memory/Grow （这个页面对返回值的描述有误）
	return previousSize
}

// 因为指令中的 offset 立即数是 uint32，而操作数栈弹出的值也是 uint32，
// 所以有效地址（uint32 + uint32）是一个 33 位的无符号整数，超出了 uint32 的范围，
// 所以这里使用 uint64 存储有效地址。
func (m *memory) Read(effective_address uint64, buf []byte) {
	if int(effective_address)+len(buf) > len(m.data) {
		panic(errors.New("out of memory boundary"))
	}
	copy(buf, m.data[effective_address:])
}

func (m *memory) Write(effective_address uint64, data []byte) {
	if int(effective_address)+len(data) > len(m.data) {
		panic(errors.New("out of memory boundary"))
	}
	copy(m.data[effective_address:], data)
}
