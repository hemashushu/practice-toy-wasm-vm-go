package interpreter

import (
	"errors"
	"wasmvm/binary"
)

// 调用函数的过程
//
// 1. 在逻辑上除了有一个 `全局变量` 表，还有一个 `局部变量` 表，函数的实参其实也是局部变量。
// 2. `调用者` 把实参准备好，存放在操作数栈顶，第一个参数位于靠近栈底，后面的参数靠近栈顶。
//
//                        局部变量表 -- | --- index 0 --- |
//                                    | 第 N 个局部变量槽  |
//                                    | 第 0 个局部变量槽  |
// 当前    | ------- 栈顶 -------- |    |................ |
// 函数    | 第 N 个参数值          | -> | 第 N 个参数值     |
// 操作数  | 第 0 个参数值          | -> | 第 0 个参数值     |
// 栈  -- | ..................... |    | --- index 0 --- |
//        |                      |
//        | 当前函数用于运算的槽位   |
//        | ------- 栈底 -------- |
//
// 3. `被调用者` 从操作数栈弹出 N 个参数，并存入局部变量表。
// 4. `被调用者` 在局部变量表开辟 N 个局部变量空槽，初始值均为 0。
// 5. `被调用者` 在操作数栈上执行当前函数的指令。
//
//                           局部变量表 -- | --- index 0 --- |
//                                       | 第 N 个局部变量槽  |
//                                       | 第 0 个局部变量槽  |
//          | ------- 栈顶 --------- |    |................ |
// 当前函数  |                       |    | 第 N 个参数值     |
// 栈帧   -- | 当前函数用于运算的槽位   |    | 第 0 个参数值     |
//          | ..................... |    | --- index 0 --- |
//          |                       |
// 调用者 -- | 上一个函数留存下来的操作数 |
// 栈帧      | ------- 栈底 --------- |
//
// 6. `被调用者` 将返回值压入操作数栈，第一个返回值先压入，后面的返回值后压入。
//
//
//        | ------- 栈顶 -------- |
//        | `被调用者` 退出后留下的  |
// 当前    | 的遗产————返回值       |
// 函数    | 第 N 个返回值          |
// 操作数  | 第 0 个返回值          |
// 栈  -- | ..................... |
//        |                      |
//        | 当前函数用于运算的槽位   |
//        | ------- 栈底 -------- |
//
// 注意：
// - 以上是函数调用的逻辑，具体的实现可能有所不同
// - 有时 `被调用者` 可能会残存一些局部变量在操作数栈上，所以在调用函数前需要
//   记录第 0 个实参的位置（地址），以便调用完目标函数后清除目标函数在操作数栈上
//   运算后的残留数据，让的返回值刚好接在调用前的栈顶（除了实参之外的位置）。
//
// 调用帧 call frame：
// 当前函数所需的数据
//
// 调用栈 call stack：
// 一连串调用帧堆起来的栈
//
// - f3  <-- 栈顶（当前调用帧，当前函数）
// - f2
// - f1
// - f0  <-- 栈底
//

// -------- call 指令
//
// call func_idx:uint32
//
// 注
// 函数索引值是包括 "导入的函数" 以及 "当前模块定义的函数（即内部函数）"，而且先计算
// 导入的函数，比如一个模块有 3 个函数导入和 2 个内部函数，则第一个内部函数的索引值为 3。

// call(vm, args) 函数仅仅创建了一个调用帧，并不会自动开始
// 执行函数当中的指令（字节码）
func call(v *vm, args interface{}) {
	idx := int(args.(uint32))
	f := v.funcs[idx]
	callFunc(v, f)
}

func callFunc(v *vm, f vmFunc) {
	// if f.goFunc != nil {
	if f.func_ != nil {
		callExternalFunc(v, f)
	} else {
		callInternalFunc(v, f)
	}
}

func callInternalFunc(v *vm, f vmFunc /* func_idx int*/) { // name: callFunction
	// funcTypeIdx := v.module.FuncSec[func_idx]
	// funcType := v.module.TypeSec[funcTypeIdx]
	// code := v.module.CodeSec[func_idx]
	// expr := code.Expr

	funcType := f.type_
	code := f.code
	expr := code.Expr

	// 创建被进入新的调用帧
	v.enterBlock(binary.Call, funcType, expr)

	// 分配局部变量空槽
	localCount := int(code.GetLocalCount())
	for i := 0; i < localCount; i++ {
		v.operandStack.pushU64(0) // 局部变量的空槽初始值为 0
	}
}

func callExternalFunc(v *vm, f vmFunc) {
	args := popArgs(v, f.type_)
	results := f.func_.Eval(args...)
	pushResults(v, f.type_, results)
}

func popArgs(v *vm, funcType binary.FuncType) []interface{} {
	paramCount := len(funcType.ParamTypes)
	args := make([]interface{}, paramCount)

	// 注：
	// 这是从模块内部函数调用外部函数的过程。
	//
	// 先弹出的参数放在参数列表的右边（大索引端）
	// 对于返回值，左边（小索引端）的数值先压入。
	//
	// 示例：
	// external function
	// extf (a,b,c) -> (x,y)
	//       ^ ^ ^      | |
	//       | | |      V V
	//
	// --- 栈顶 ---    --- 栈顶 ---
	// - c
	// - b            - y
	// - a            - x
	// - ...          - ...
	// --- 栈底 ---    --- 栈顶 ---

	for i := paramCount - 1; i >= 0; i-- {
		args[i] = wrapU64(funcType.ParamTypes[i], v.operandStack.popU64())
	}
	return args
}

func pushResults(v *vm, ft binary.FuncType, results []interface{}) {
	if len(ft.ResultTypes) != len(results) {
		panic(errors.New("incorrect length of return values"))
	}
	for _, result := range results {
		v.operandStack.pushU64(unwrapU64(ft.ResultTypes[0], result))
	}
}

// -------- call_indirect 间接函数调用
//
// call_indirect type_idx:uint32 table_idx:uint32
//
// 其中 table_idx 的值目前只能是 0
//
func callIndirect(v *vm, args interface{}) {
	i := v.operandStack.popU32() // 读取目标表项的索引
	if i >= v.table.Size() {
		panic(errors.New("out of element range"))
	}

	f := v.table.GetElem(i)

	typeIdx := args.(uint32)
	funcType := v.module.TypeSec[typeIdx]

	// todo::
	// 这里需要检查函数类型（函数的签名）是否匹配
	// if f.type_.GetSignature() != funcType.GetSignature() {
	// 	panic(errors.New("function type mismatch in indirect call"))
	// }
	// call_indirect 指令的 type_idx 参数用于防止调用错了函数

	// callFunc(v, f)

	// 检查是否同一个模块里的内部函数
	if vf, ok := f.(vmFunc); ok {
		if vf.vm == v && vf.func_ == nil {
			callInternalFunc(v, vf)
			return
		}
	}

	fcArgs := popArgs(v, funcType)
	results := f.Eval(fcArgs...)
	pushResults(v, funcType, results)
}
