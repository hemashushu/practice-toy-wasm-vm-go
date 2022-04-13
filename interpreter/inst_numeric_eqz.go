package interpreter

// ======== 数值指令
//
// -------- 等零测试指令
//
// 从栈顶弹出一个操作数，判断是否为 0，
// 如果为 0 则压入 1（int32）， 否则压入 0（int32）
//
// i32.eqz
// i64.eqz

func i32Eqz(v *vm, _ interface{}) {
	v.operandStack.pushBool(v.operandStack.popS32() == 0)
}

func i64Eqz(v *vm, _ interface{}) {
	v.operandStack.pushBool(v.operandStack.popS64() == 0)
}
