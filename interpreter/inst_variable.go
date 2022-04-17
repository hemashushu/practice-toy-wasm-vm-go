package interpreter

// ======== 变量指令
//
// 读写局部/全局变量
//
// local.get local_idx:uint32	;; 读取指定索引的局部变量的值，压入操作数栈
// local.set local_idx:uint32	;; 从操作数栈弹出一个数，写入到指定索引的局部变量；弹出的数的类型必须跟局部变量的一致
// local.tee local_idx:uint32	;; 读取栈顶的值，写入到指定索引的局部变量
//
// global.get global_idx:uint32	;; 读取指定索引的全局变量的值，压入操作数栈
// global.set global_idx:uint32	;; 从操作数栈弹出一个数，写入到指定索引的全局变量；弹出的数的类型必须跟全局变量的一致

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
	val := v.operandStack.peekValue()
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
