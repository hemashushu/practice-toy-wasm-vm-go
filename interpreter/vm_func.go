package interpreter

import (
	"errors"
	"wasmvm/binary"
	"wasmvm/instance"
)

// 使用数组和传递参数和返回值
// 因为 wasm 只支持 4 种数据类型，所以数组的元素的数据类型
// 也可以使用 uint64。

// type GoFunc = func(args []WasmVal) []WasmVal
// type WasmVal = interface{}

type vmFunc struct {
	type_ binary.FuncType // name: func_type

	code binary.Code // code 和 goFunc 二选一
	vm   *vm

	// goFunc GoFunc          // 本地函数（native function）
	func_ instance.Function // 外部函数，即从别的模块导入的函数
}

func newExternalFunc(
	funcType binary.FuncType,
	//f GoFunc
	f instance.Function) vmFunc {
	return vmFunc{
		type_: funcType,
		// goFunc: f,
		func_: f,
	}
}
func newInternalFunc(
	v *vm,
	funcType binary.FuncType,
	code binary.Code) vmFunc {
	return vmFunc{
		type_: funcType,
		code:  code,
		vm:    v,
	}
}

func (f vmFunc) Type() binary.FuncType {
	return f.type_
}

// name: CallFromHost
func (f vmFunc) Eval(args ...instance.WasmVal) []instance.WasmVal {
	if f.func_ != nil {
		// 外部函数
		return f.func_.Eval(args...)
	} else {
		// 模块内部函数
		return f.eval(args)
	}
}

// 从 vm 外部调用模块内部的函数
func (f vmFunc) eval(args []interface{}) []interface{} {
	pushArgs(f.vm, f.type_, args)
	callFunc(f.vm, f)
	if f.func_ == nil {
		f.vm.loop()
	}
	return popResults(f.vm, f.type_)
}

func pushArgs(v *vm, ft binary.FuncType, args []interface{}) {

	// 注：
	// 这是从外部函数调用模块内部函数的过程。
	//
	// 参数列表左边（小索引端）的实参先压入
	// 对于返回结果，先弹出的数值放置在结果数组的右边（大索引端）
	//
	// 示例：
	// (a,b,c) -> (x,y)
	//  | | |      ^ ^
	//  V V V      | |
	//
	// internal function
	//
	// --- 栈顶 ---    --- 栈顶 ---
	// - c
	// - b            - y
	// - a            - x
	// - ...          - ...
	// --- 栈底 ---    --- 栈顶 ---

	if len(ft.ParamTypes) != len(args) {
		panic(errors.New("incorrect length of arguments"))
	}
	for i, vt := range ft.ParamTypes {
		v.operandStack.pushU64(unwrapU64(vt, args[i]))
	}
}
func popResults(v *vm, ft binary.FuncType) []interface{} {
	results := make([]interface{}, len(ft.ResultTypes))
	for n := len(ft.ResultTypes) - 1; n >= 0; n-- {
		results[n] = wrapU64(ft.ResultTypes[n], v.operandStack.popU64())
	}
	return results
}
