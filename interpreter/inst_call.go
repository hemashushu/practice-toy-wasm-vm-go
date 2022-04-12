package interpreter

import "fmt"

// hack!
func call(v *vm, args interface{}) {
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
