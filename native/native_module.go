package native

import (
	"wasmvm/binary"
	"wasmvm/instance"
)

type nativeModule struct {
	exported map[string]interface{}
}

func NewNativeModule() nativeModule {
	return nativeModule{exported: map[string]interface{}{}}
}

// func (m nativeModule) Register(name string, value interface{}) {
// 	m.exported[name] = value
// }

func (m nativeModule) RegisterFunc(name string, paramTypes []binary.ValType, resultTypes []binary.ValType, func_ GoFunc) {
	funcType := binary.FuncType{
		ParamTypes:  paramTypes,
		ResultTypes: resultTypes,
		Tag:         binary.FtTag}
	m.exported[name] = nativeFunction{funcType: funcType, func_: func_}
}

// 实现接口 Module 的方法

func (m nativeModule) GetMember(name string) interface{} {
	return m.exported[name]
}

func (m nativeModule) EvalFunc(name string, args ...instance.WasmVal) []instance.WasmVal {
	return m.exported[name].(instance.Function).Eval(args...)
}

func (m nativeModule) GetGlobalVal(name string) instance.WasmVal {
	return m.exported[name].(instance.Global).Get()
}

func (m nativeModule) SetGlobalVal(name string, value instance.WasmVal) {
	m.exported[name].(instance.Global).Set(value)
}

type GoFunc = func(args []instance.WasmVal) []instance.WasmVal

type nativeFunction struct {
	funcType binary.FuncType
	func_    GoFunc
}

func (f nativeFunction) Type() binary.FuncType {
	return f.funcType
}

func (f nativeFunction) Eval(args ...instance.WasmVal) []instance.WasmVal {
	return f.func_(args)
}
