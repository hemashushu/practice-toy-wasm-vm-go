package interpreter

// ======== 数值指令
//
// -------- 比较测试指令
//
// 比较测试包括 `相等测试` 以及 `大小比较测试`
//
// 从栈顶弹出 2 个操作数，然后把比较结果（int32，相当于 boolean）压入栈
// 注意先弹出的作为 RHS，后弹出的作为 LHS。

// i32

func i32Eq(v *vm, _ interface{}) {
	// 相等比较既可以作为 unsigned 数字取出，也可以作为 signed 数字取出
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushBool(lhs == rhs)
}

func i32Ne(v *vm, _ interface{}) {
	// 相等比较既可以作为 unsigned 数字取出，也可以作为 signed 数字取出
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushBool(lhs != rhs)
}

func i32LtS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS32(), v.operandStack.popS32()
	v.operandStack.pushBool(lhs < rhs)
}

func i32LtU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushBool(lhs < rhs)
}

func i32GtS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS32(), v.operandStack.popS32()
	v.operandStack.pushBool(lhs > rhs)
}

func i32GtU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushBool(lhs > rhs)
}

func i32LeS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS32(), v.operandStack.popS32()
	v.operandStack.pushBool(lhs <= rhs)
}

func i32LeU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushBool(lhs <= rhs)
}

func i32GeS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS32(), v.operandStack.popS32()
	v.operandStack.pushBool(lhs >= rhs)
}

func i32GeU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushBool(lhs >= rhs)
}

// i64

func i64Eq(v *vm, _ interface{}) {
	// 相等比较既可以作为 unsigned 数字取出，也可以作为 signed 数字取出
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushBool(lhs == rhs)
}

func i64Ne(v *vm, _ interface{}) {
	// 相等比较既可以作为 unsigned 数字取出，也可以作为 signed 数字取出
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushBool(lhs != rhs)
}

func i64LtS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS64(), v.operandStack.popS64()
	v.operandStack.pushBool(lhs < rhs)
}

func i64LtU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushBool(lhs < rhs)
}

func i64GtS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS64(), v.operandStack.popS64()
	v.operandStack.pushBool(lhs > rhs)
}

func i64GtU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushBool(lhs > rhs)
}

func i64LeS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS64(), v.operandStack.popS64()
	v.operandStack.pushBool(lhs <= rhs)
}

func i64LeU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushBool(lhs <= rhs)
}

func i64GeS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS64(), v.operandStack.popS64()
	v.operandStack.pushBool(lhs >= rhs)
}

func i64GeU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushBool(lhs >= rhs)
}

// f32

func f32Eq(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushBool(lhs == rhs)
}

func f32Ne(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushBool(lhs != rhs)
}

func f32Lt(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushBool(lhs < rhs)
}

func f32Gt(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushBool(lhs > rhs)
}

func f32Le(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushBool(lhs <= rhs)
}

func f32Ge(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushBool(lhs >= rhs)
}

// f64

func f64Eq(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushBool(lhs == rhs)
}

func f64Ne(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushBool(lhs != rhs)
}

func f64Lt(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushBool(lhs < rhs)
}

func f64Gt(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushBool(lhs > rhs)
}

func f64Le(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushBool(lhs <= rhs)
}

func f64Ge(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushBool(lhs >= rhs)
}
