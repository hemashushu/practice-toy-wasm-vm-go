package interpreter

import (
	"errors"
	"wasmvm/binary"
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

func (g *globalVar) GetAsU64() uint64 {
	return g.val
}
func (g *globalVar) SetAsU64(val uint64) {
	if g.type_.Mut != 1 {
		panic(errors.New("immutable global"))
	}
	g.val = val
}
