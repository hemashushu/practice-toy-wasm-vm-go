package interpreter

// ======== 数值指令
//
// -------- 常量指令
// 将指令当中的立即数压入栈
//
// i32.const
// i64.const
// f32.const
// f64.const

func i32Const(v *vm, args interface{}) {
	v.operandStack.pushS32(args.(int32))
}

func i64Const(v *vm, args interface{}) {
	v.operandStack.pushS64(args.(int64))
}

func f32Const(v *vm, args interface{}) {
	v.operandStack.pushF32(args.(float32))
}
func f64Const(v *vm, args interface{}) {
	v.operandStack.pushF64(args.(float64))
}
