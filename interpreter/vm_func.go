package interpreter

import "wasmvm/binary"

// 使用数组和传递参数和返回值
// 因为 wasm 只支持 4 种数据类型，所以数组的元素的数据类型
// 也可以使用 uint64。
type GoFunc = func(args []WasmVal) []WasmVal
type WasmVal = interface{}

type vmFunc struct {
	type_  binary.FuncType // name: func_type
	code   binary.Code     // code 和 goFunc 二选一
	goFunc GoFunc          // 本地函数（native function）
}

func newExternalFunc(funcType binary.FuncType, f GoFunc) vmFunc {
	return vmFunc{
		type_:  funcType,
		goFunc: f,
	}
}
func newInternalFunc(funcType binary.FuncType, code binary.Code) vmFunc {
	return vmFunc{
		type_: funcType,
		code:  code,
	}
}

func (f vmFunc) Type() binary.FuncType {
	return f.type_
}
