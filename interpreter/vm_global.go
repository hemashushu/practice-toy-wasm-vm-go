package interpreter

import (
	"errors"
	"wasmvm/binary"
	"wasmvm/instance"
)

type globalVar struct {
	// GlobalType 记录着：`数据类型` 以及 `是否可变` 两样信息
	type_ binary.GlobalType

	// 数值
	val uint64
}

func newGlobal(gt binary.GlobalType, val uint64) *globalVar {
	return &globalVar{type_: gt, val: val}
}

func (g *globalVar) Type() binary.GlobalType {
	return g.type_
}

func (g *globalVar) GetAsU64() uint64 { // 内部使用，name: GetRaw()
	return g.val
}

func (g *globalVar) SetAsU64(val uint64) { // 内部使用，name: SetRaw(...)
	if g.type_.Mut != 1 {
		panic(errors.New("immutable global"))
	}
	g.val = val
}

func (g *globalVar) Get() instance.WasmVal {
	return wrapU64(g.type_.ValType, g.val)
}

func (g *globalVar) Set(val instance.WasmVal) {
	g.val = unwrapU64(g.type_.ValType, val)
}
