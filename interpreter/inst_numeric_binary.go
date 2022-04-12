package interpreter

import (
	"errors"
	"math"
	"math/bits"
)

// ======== 数值指令
//
// -------- 二元运算指令
//
// 从栈顶弹出 2 个操作数，经过运算后压入栈
// 先弹出的是 rhs，后弹出的是 lhs
//
// 注：
// - 整数的旋转指令 rotl 和 rotr 使用了 go 标准库实现
// - 浮点数的 min, max 使用了 go 标准库实现
// - 浮点数的 copysign 的作用：应用位于栈顶的操作数的符号到位于第二位的操作数
// - 移位操作可以先对第二个操作数进行一次求余运算，以防止错误

// i32

func i32Add(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs + rhs)
}

func i32Sub(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs - rhs)
}

func i32Mul(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs * rhs)
}

func i32DivS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS32(), v.operandStack.popS32()
	if lhs == math.MinInt32 && rhs == -1 {
		panic(errors.New("integer overflow"))
	}
	v.operandStack.pushS32(lhs / rhs)
}

func i32DivU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs / rhs)
}

func i32RemS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS32(), v.operandStack.popS32()
	v.operandStack.pushS32(lhs % rhs)
}

func i32RemU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs % rhs)
}

func i32And(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs & rhs)
}

func i32Or(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs | rhs)
}

func i32Xor(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs ^ rhs)
}

func i32Shl(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs << (rhs % 32))
}

func i32ShrS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popS32()
	v.operandStack.pushS32(lhs >> (rhs % 32))
}

func i32ShrU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(lhs >> (rhs % 32))
}

func i32Rotl(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(bits.RotateLeft32(lhs, int(rhs)))
}

func i32Rotr(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU32(), v.operandStack.popU32()
	v.operandStack.pushU32(bits.RotateLeft32(lhs, -int(rhs)))
}

// i64

func i64Add(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs + rhs)
}

func i64Sub(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs - rhs)
}

func i64Mul(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs * rhs)
}

func i64DivS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS64(), v.operandStack.popS64()
	if lhs == math.MinInt64 && rhs == -1 {
		panic(errors.New("integer overflow"))
	}
	v.operandStack.pushS64(lhs / rhs)
}

func i64DivU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs / rhs)
}

func i64RemS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popS64(), v.operandStack.popS64()
	v.operandStack.pushS64(lhs % rhs)
}

func i64RemU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs % rhs)
}

func i64And(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs & rhs)
}

func i64Or(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs | rhs)
}

func i64Xor(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs ^ rhs)
}

func i64Shl(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs << (rhs % 64))
}

func i64ShrS(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popS64()
	v.operandStack.pushS64(lhs >> (rhs % 64))
}

func i64ShrU(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(lhs >> (rhs % 64))
}

func i64Rotl(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(bits.RotateLeft64(lhs, int(rhs)))
}

func i64Rotr(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popU64(), v.operandStack.popU64()
	v.operandStack.pushU64(bits.RotateLeft64(lhs, -int(rhs)))
}

// f32

func f32Add(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushF32(lhs + rhs)
}

func f32Sub(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushF32(lhs - rhs)
}

func f32Mul(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushF32(lhs * rhs)
}

func f32Div(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushF32(lhs / rhs)
}

func f32Min(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	leftNaN := math.IsNaN(float64(lhs))
	rightNaN := math.IsNaN(float64(rhs))
	if leftNaN && !rightNaN {
		v.operandStack.pushF32(lhs)
		return
	} else if rightNaN && !leftNaN {
		v.operandStack.pushF32(rhs)
		return
	}
	v.operandStack.pushF32(float32(math.Min(float64(lhs), float64(rhs))))
}

func f32Max(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF32(), v.operandStack.popF32()
	leftNaN := math.IsNaN(float64(lhs))
	rightNaN := math.IsNaN(float64(rhs))
	if leftNaN && !rightNaN {
		v.operandStack.pushF32(lhs)
		return
	} else if rightNaN && !leftNaN {
		v.operandStack.pushF32(rhs)
		return
	}
	v.operandStack.pushF32(float32(math.Max(float64(lhs), float64(rhs))))
}

func f32CopySign(v *vm, _ interface{}) {
	from, to := v.operandStack.popF32(), v.operandStack.popF32()
	v.operandStack.pushF32(float32(math.Copysign(float64(to), float64(from))))
}

// f64

func f64Add(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushF64(lhs + rhs)
}

func f64Sub(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushF64(lhs - rhs)
}

func f64Mul(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushF64(lhs * rhs)
}

func f64Div(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushF64(lhs / rhs)
}

func f64Min(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	leftNaN := math.IsNaN(lhs)
	rightNaN := math.IsNaN(rhs)
	if leftNaN && !rightNaN {
		v.operandStack.pushF64(lhs)
		return
	} else if rightNaN && !leftNaN {
		v.operandStack.pushF64(rhs)
		return
	}
	v.operandStack.pushF64(math.Min(lhs, rhs))
}

func f64Max(v *vm, _ interface{}) {
	rhs, lhs := v.operandStack.popF64(), v.operandStack.popF64()
	leftNaN := math.IsNaN(lhs)
	rightNaN := math.IsNaN(rhs)
	if leftNaN && !rightNaN {
		v.operandStack.pushF64(lhs)
		return
	} else if rightNaN && !leftNaN {
		v.operandStack.pushF64(rhs)
		return
	}
	v.operandStack.pushF64(math.Max(lhs, rhs))
}

func f64CopySign(v *vm, _ interface{}) {
	from, to := v.operandStack.popF64(), v.operandStack.popF64()
	v.operandStack.pushF64(math.Copysign(to, from))
}
