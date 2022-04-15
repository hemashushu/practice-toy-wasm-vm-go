package interpreter

// ======== 操作数（参数 parametric）指令
//
// 用于修改操作数栈元素的指令，包括 drop 和 select

// ### drop
//
// 弹出栈顶的一个操作数并扔掉

func drop(v *vm, _ interface{}) {
	v.operandStack.popU64()
}

// ### select
//
// 从栈弹出 3 个操作数，根据栈顶操作数（int32）是否为零，
// 来决定是压入第 2 个操作数（consequent）或者第 3 个操作数（alternate）
//
// - testing     <-- 栈顶
// - consequent  <-- 为 0 时压入这个
// - alternate   <-- 为 1 时压入这个
//
// 示例：
// (i32.const 11)  <-- alternate
// (i32.const 10)  <-- consequent
// (i32.const 0)   <-- testing
// (select)        <-- 最终压入 10

// 其中：
// 栈顶元素（第一个操作数）必须是 int32，
// 第二个和第三个操作数的类型必须相同

func select_(v *vm, _ interface{}) {
	testing, consequent, alternate := v.operandStack.popU32(), v.operandStack.popU64(), v.operandStack.popU64()
	if testing == 0 {
		v.operandStack.pushU64(consequent)
	} else {
		v.operandStack.pushU64(alternate)
	}
}
