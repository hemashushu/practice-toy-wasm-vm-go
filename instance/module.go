package instance

import "wasmvm/binary"

// 模块实例
type Module interface {
	// 根据名称获取导出项
	// 导出项只能是：函数、表、内存、全局变量
	GetMember(name string) interface{} // name: getExportItem

	// 辅助函数
	InvokeFunc(name string, args ...WasmVal) []WasmVal
	GetGlobalVal(name string) WasmVal
	SetGlobalVal(name string, value WasmVal)
}

type WasmVal = interface{}

// 导出项 -- 函数
type Function interface {
	Type() binary.FuncType
	Eval(args ...WasmVal) []WasmVal
}

// 导出项 -- 表
type Table interface {
	Type() binary.TableType
	Size() uint32
	Grow(increaseNumber uint32)
	GetElem(idx uint32) Function
	SetElem(idx uint32, elem Function)
}

// 导出项 -- 内存
type Memory interface {
	Type() binary.MemType
	Size() uint32                      // 注，是页面数量，而不是字节数
	Grow(increaseNumber uint32) uint32 // increaseNumber: 需增加的页面数，失败时会返回被转为 uint32 的 -1
	Read(offset uint64, buf []byte)
	Write(offset uint64, buf []byte)
}

type Global interface {
	Type() binary.GlobalType
	GetAsU64() uint64      // 内部使用，name: GetRaw()
	SetAsU64(value uint64) // 内部使用，name: SetRaw(...)
	Get() WasmVal
	Set(value WasmVal)
}
