package interpreter

import (
	"errors"
	"math"
)

// ======== 数值指令
//
// -------- 类型转换指令
//
// ### 整数截断
//
// i32.wrap_i64
//
// 将 64 位的整数直接截断为 32 位（即只保留低端信息）

func i32WrapI64(v *vm, _ interface{}) {
	v.operandStack.pushU32(uint32(v.operandStack.popU64()))
}

// ### 整数提升
//
// 将位宽较窄的整数提升为位宽较广的整数，比如将 32 位整数提升为 64 位
//
// 源 i32，目标 i32
// i32.extend8_s
// i32.extend16_s
//
// 源 i32，目标 i64
// i64.extend_i32_s
// i64.extend_i32_u
//
// 源 i64，目标 i64
// i64.extend8_s
// i64.extend16_s
// i64.extend32_s

func i32Extend8S(v *vm, _ interface{}) {
	v.operandStack.pushS32(int32(int8(v.operandStack.popS32())))
}

func i32Extend16S(v *vm, _ interface{}) {
	v.operandStack.pushS32(int32(int16(v.operandStack.popS32())))
}

func i64ExtendI32S(v *vm, _ interface{}) {
	v.operandStack.pushS64(int64(v.operandStack.popS32()))
}

func i64ExtendI32U(v *vm, _ interface{}) {
	v.operandStack.pushU64(uint64(v.operandStack.popU32()))
}

func i64Extend8S(v *vm, _ interface{}) {
	v.operandStack.pushS64(int64(int8(v.operandStack.popS64())))
}

func i64Extend16S(v *vm, _ interface{}) {
	v.operandStack.pushS64(int64(int16(v.operandStack.popS64())))
}

func i64Extend32S(v *vm, _ interface{}) {
	v.operandStack.pushS64(int64(int32(v.operandStack.popS64())))
}

// ### 浮点数转整数（截断运算）
//
// 把浮点数截断为整数
//
// 源 f32，目标 i32
// i32.trunc_f32_s
// i32.trunc_f32_u
//
// 源 f32，目标 i64
// i64.trunc_f32_s
// i64.trunc_f32_u
//
// 源 f64，目标 i32
// i32.trunc_f64_s
// i32.trunc_f64_u
//
// 源 f64，目标 i64
// i64.trunc_f64_s
// i64.trunc_f64_u

func i32TruncF32S(v *vm, _ interface{}) {
	f := math.Trunc(float64(v.operandStack.popF32()))
	if f > math.MaxInt32 || f < math.MinInt32 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushS32(int32(f))
}

func i32TruncF32U(v *vm, _ interface{}) {
	f := math.Trunc(float64(v.operandStack.popF32()))
	if f > math.MaxUint32 || f < 0 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushU32(uint32(f))
}

func i64TruncF32S(v *vm, _ interface{}) {
	f := math.Trunc(float64(v.operandStack.popF32()))
	if f >= math.MaxInt64 || f < math.MinInt64 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushS64(int64(f))
}

func i64TruncF32U(v *vm, _ interface{}) {
	f := math.Trunc(float64(v.operandStack.popF32()))
	if f >= math.MaxUint64 || f < 0 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushU64(uint64(f))
}

func i32TruncF64S(v *vm, _ interface{}) {
	f := math.Trunc(v.operandStack.popF64())
	if f > math.MaxInt32 || f < math.MinInt32 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushS32(int32(f))
}

func i32TruncF64U(v *vm, _ interface{}) {
	f := math.Trunc(v.operandStack.popF64())
	if f > math.MaxUint32 || f < 0 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushU32(uint32(f))
}

func i64TruncF64S(v *vm, _ interface{}) {
	f := math.Trunc(v.operandStack.popF64())
	if f >= math.MaxInt64 || f < math.MinInt64 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushS64(int64(f))
}

func i64TruncF64U(v *vm, _ interface{}) {
	f := math.Trunc(v.operandStack.popF64())
	if f >= math.MaxUint64 || f < 0 {
		panic(errors.New("integer overflow"))
	}
	if math.IsNaN(f) {
		panic(errors.New("invalid conversion to integer"))
	}
	v.operandStack.pushU64(uint64(f))
}

// ### 饱和截断
//
// 跟一般截断不同的是：
// - 将 NaN 转为 0
// - 将正/负无穷转为整数最大/最小值

func truncSat(v *vm, args interface{}) {
	switch args.(byte) {
	case 0: // i32.trunc_sat_f32_s
		val := truncSatS(float64(v.operandStack.popF32()), 32)
		v.operandStack.pushS32(int32(val))
	case 1: // i32.trunc_sat_f32_u
		val := truncSatU(float64(v.operandStack.popF32()), 32)
		v.operandStack.pushU32(uint32(val))
	case 2: // i32.trunc_sat_f64_s
		val := truncSatS(v.operandStack.popF64(), 32)
		v.operandStack.pushS32(int32(val))
	case 3: // i32.trunc_sat_f64_u
		val := truncSatU(v.operandStack.popF64(), 32)
		v.operandStack.pushU32(uint32(val))
	case 4: // i64.trunc_sat_f32_s
		val := truncSatS(float64(v.operandStack.popF32()), 64)
		v.operandStack.pushS64(val)
	case 5: // i64.trunc_sat_f32_u
		val := truncSatU(float64(v.operandStack.popF32()), 64)
		v.operandStack.pushU64(val)
	case 6: // i64.trunc_sat_f64_s
		val := truncSatS(v.operandStack.popF64(), 64)
		v.operandStack.pushS64(val)
	case 7: // i64.trunc_sat_f64_u
		val := truncSatU(v.operandStack.popF64(), 64)
		v.operandStack.pushU64(val)
	default:
		panic(errors.New("unreachable"))
	}
}

func truncSatU(z float64, n int) uint64 {
	if math.IsNaN(z) {
		return 0
	}
	if math.IsInf(z, -1) {
		return 0
	}
	max := (uint64(1) << n) - 1
	if math.IsInf(z, 1) {
		return max
	}
	if x := math.Trunc(z); x < 0 {
		return 0
	} else if x >= float64(max) {
		return max
	} else {
		return uint64(x)
	}
}

func truncSatS(z float64, n int) int64 {
	if math.IsNaN(z) {
		return 0
	}
	min := -(int64(1) << (n - 1))
	max := (int64(1) << (n - 1)) - 1
	if math.IsInf(z, -1) {
		return min
	}
	if math.IsInf(z, 1) {
		return max
	}
	if x := math.Trunc(z); x < float64(min) {
		return min
	} else if x >= float64(max) {
		return max
	} else {
		return int64(x)
	}
}

// ### 整数转浮点数（转换运算）
//
// 源 i32，目标 f32
// f32.convert_i32_s
// f32.convert_i32_u
//
// 源 i32，目标 f64
// f64.convert_i32_s
// f64.convert_i32_u
//
// 源 i64，目标 f32
// f32.convert_i64_s
// f32.convert_i64_u
//
// 源 i64，目标 f64
// f64.convert_i64_s
// f64.convert_i64_u

func f32ConvertI32S(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(v.operandStack.popS32()))
}

func f32ConvertI32U(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(v.operandStack.popU32()))
}

func f64ConvertI32S(v *vm, _ interface{}) {
	v.operandStack.pushF64(float64(v.operandStack.popS32()))
}

func f64ConvertI32U(v *vm, _ interface{}) {
	v.operandStack.pushF64(float64(v.operandStack.popU32()))
}

func f32ConvertI64S(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(v.operandStack.popS64()))
}

func f32ConvertI64U(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(v.operandStack.popU64()))
}

func f64ConvertI64S(v *vm, _ interface{}) {
	v.operandStack.pushF64(float64(v.operandStack.popS64()))
}

func f64ConvertI64U(v *vm, _ interface{}) {
	v.operandStack.pushF64(float64(v.operandStack.popU64()))
}

// ### 浮点数精度调整
//
// f32.demote_f64_s
// f64.promote_f32

func f32DemoteF64(v *vm, _ interface{}) {
	v.operandStack.pushF32(float32(v.operandStack.popF64()))
}

func f64PromoteF32(v *vm, _ interface{}) {
	v.operandStack.pushF64(float64(v.operandStack.popF32()))
}

// ### 比特位重新解释
//
// 不改变操作数的比特位，仅重新解释成其他类型

func i32ReinterpretF32(v *vm, _ interface{}) {
	// 当前的操作数栈实现已经统一转换为 uint64，所以这里无需操作
}

func i64ReinterpretF64(v *vm, _ interface{}) {
	// 当前的操作数栈实现已经统一转换为 uint64，所以这里无需操作
}

func f32ReinterpretI32(v *vm, _ interface{}) {
	// 当前的操作数栈实现已经统一转换为 uint64，所以这里无需操作
}

func f64ReinterpretI64(v *vm, _ interface{}) {
	// 当前的操作数栈实现已经统一转换为 uint64，所以这里无需操作
}
