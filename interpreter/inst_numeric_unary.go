package interpreter

import (
	"math"
	"math/bits"
)

// ======== 数值指令
//
// -------- 一元运算指令
//
// 从栈顶弹出一个操作数，经过运算后压入栈
//
// i32.clz:		count leading zeros 统计前缀（高端位）比特值为 0 的数量 0x67
// i32.ctz:		count trailing zeros 统计后缀（低端位）比特值为 0 的数量 0x68
// i32.popcnt:  population count 统计所有位当中，比特值为 1 的数量 0x69
//
// 如 8'b00001100, clz == 4, ctz == 2, popcnt == 2
//
// 浮点数的一元运算，返回值仍然是浮点数
//
// f32.abs:		绝对值
// f32.neg:		取反
// f32.ceil:	向上取整（x 数轴向右方向取整）
// f32.floor:   向下取整（x 数轴向左方向取整）
// f32.trunc:   直接裁掉小数部分，取整
// f32.nearest: 就近取整（4 舍 6 入，5 奇进偶不进）
// f32.sqrt:	平方根
//
// 注：
// - 只有 neg 可以直接映射到 go 的一元运算，其他的使用了标准库实现
//
// 关于 nearest 函数
// https://en.wikipedia.org/wiki/Rounding
// https://developer.mozilla.org/en-US/docs/WebAssembly/Reference/Numeric/Nearest

// i32

func i32Clz(v *vm, _ interface{}) {
	v.operandStack.pushU32(uint32(bits.LeadingZeros32(v.operandStack.popU32())))
}

func i32Ctz(v *vm, _ interface{}) {
	v.operandStack.pushU32(uint32(bits.TrailingZeros32(v.operandStack.popU32())))
}

func i32PopCnt(v *vm, _ interface{}) {
	v.operandStack.pushU32(uint32(bits.OnesCount32(v.operandStack.popU32())))
}

// i64

func i64Clz(v *vm, _ interface{}) {
	v.operandStack.pushU64(uint64(bits.LeadingZeros64(v.operandStack.popU64())))
}

func i64Ctz(v *vm, _ interface{}) {
	v.operandStack.pushU64(uint64(bits.TrailingZeros64(v.operandStack.popU64())))
}

func i64PopCnt(v *vm, _ interface{}) {
	v.operandStack.pushU64(uint64(bits.OnesCount64(v.operandStack.popU64())))
}

// f32

func f32Abs(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(math.Abs(float64(v.operandStack.popF32()))))
}

func f32Neg(v *vm, _ interface{}) {
	v.operandStack.pushF32(-v.operandStack.popF32())
}

func f32Ceil(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(math.Ceil(float64(v.operandStack.popF32()))))
}

func f32Floor(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(math.Floor(float64(v.operandStack.popF32()))))
}

func f32Trunc(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(math.Trunc(float64(v.operandStack.popF32()))))
}

func f32Nearest(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(math.RoundToEven(float64(v.operandStack.popF32()))))
}

func f32Sqrt(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(math.Sqrt(float64(v.operandStack.popF32()))))
}

// f64

func f64Abs(v *vm, _ interface{}) {
	v.operandStack.pushF64(math.Abs(v.operandStack.popF64()))
}

func f64Neg(v *vm, _ interface{}) {
	v.operandStack.pushF64(-v.operandStack.popF64())
}

func f64Ceil(v *vm, _ interface{}) {
	v.operandStack.pushF64(math.Ceil(v.operandStack.popF64()))
}

func f64Floor(v *vm, _ interface{}) {
	v.operandStack.pushF64(math.Floor(v.operandStack.popF64()))
}

func f64Trunc(v *vm, _ interface{}) {
	v.operandStack.pushF64(math.Trunc(v.operandStack.popF64()))
}

func f64Nearest(v *vm, _ interface{}) {
	v.operandStack.pushF64(math.RoundToEven(v.operandStack.popF64()))
}

func f64Sqrt(v *vm, _ interface{}) {
	v.operandStack.pushF64(math.Sqrt(v.operandStack.popF64()))
}
