package interpreter

import (
	"fmt"
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

func call(v *vm, args interface{}) {
	idx := int(args.(uint32))
	importedFuncCount := len(v.module.ImportSec)
	if idx < importedFuncCount {
		callAssertFunc(v, args)
	} else {
		callInteralFunc(v, idx-importedFuncCount)
	}
}

func callInteralFunc(v *vm, func_idx int) { // name: callFunction
	funcTypeIdx := v.module.FuncSec[func_idx]
	funcType := v.module.TypeSec[funcTypeIdx]
	code := v.module.CodeSec[func_idx]

	// 创建被进入新的调用帧
	v.enterBlock(binary.Call, funcType, code.Expr)

	// 分配局部变量空槽
	localCount := int(code.GetLocalCount())
	for i := 0; i < localCount; i++ {
		v.operandStack.pushU64(0) // 局部变量的空槽初始值为 0
	}
}

func callAssertFunc(v *vm, args interface{}) {
	idx := args.(uint32)
	switch v.module.ImportSec[idx].Name {
	case "assert_true":
		assertEq(v.operandStack.popBool(), true)
	case "assert_false":
		assertEq(v.operandStack.popBool(), false)
	case "assert_eq_i32":
		assertEq(v.operandStack.popU32(), v.operandStack.popU32())
	case "assert_eq_i64":
		assertEq(v.operandStack.popU64(), v.operandStack.popU64())
	case "assert_eq_f32":
		assertEq(v.operandStack.popF32(), v.operandStack.popF32())
	case "assert_eq_f64":
		assertEq(v.operandStack.popF64(), v.operandStack.popF64())
	default:
		panic("TODO")
	}
}

func assertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Errorf("%v != %v", a, b))
	}
}
