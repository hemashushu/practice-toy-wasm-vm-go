package interpreter

// ======== 变量指令
//
// 读写局部/全局变量
//
// local.get local_idx:uint32	;; 读取指定索引的局部变量的值，压入操作数栈

func localGet(v *vm, args interface{}) {
	idx := args.(uint32)
	val := v.operandStack.getOperand(v.local0Idx + idx)
	v.operandStack.pushU64(val)
}

func localSet(v *vm, args interface{}) {
	idx := args.(uint32)
	val := v.operandStack.popU64()
	v.operandStack.setOperand(v.local0Idx+idx, val)
}

func localTee(v *vm, args interface{}) {
	idx := args.(uint32)
	// val := vm.operandStack.popU64()
	// vm.operandStack.pushU64(val)
	val := v.operandStack.peekU64()
	v.operandStack.setOperand(v.local0Idx+idx, val)
}

func globalGet(v *vm, args interface{}) {
	idx := args.(uint32)
	val := v.globals[idx].GetAsU64()
	v.operandStack.pushU64(val)
}

func globalSet(v *vm, args interface{}) {
	idx := args.(uint32)
	val := v.operandStack.popU64()
	v.globals[idx].SetAsU64(val)
}
