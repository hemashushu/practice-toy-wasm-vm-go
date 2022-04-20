package native

import (
	"fmt"
	"wasmvm/binary"
	"wasmvm/instance"
)

func NewEnvModule() instance.Module {
	env := NewNativeModule()

	env.RegisterFunc(
		"add_i32",
		[]byte{binary.ValTypeI32, binary.ValTypeI32},
		[]byte{binary.ValTypeI32},
		addI32)

	env.RegisterFunc(
		"print_char",
		[]byte{binary.ValTypeI32},
		[]byte{},
		printChar)

	env.RegisterFunc(
		"print_int",
		[]byte{binary.ValTypeI32},
		[]byte{},
		printInt)

	return env
}

// 用于测试
func printChar(args []interface{}) []interface{} {
	fmt.Printf("%c", args[0].(int32))
	return nil
}

// 用于测试
func printInt(args []interface{}) []interface{} {
	fmt.Printf("%d", args[0].(int32))
	return nil
}

// 用于 native function 单元测试
func addI32(args []interface{}) []interface{} {
	a := args[0].(int32)
	b := args[1].(int32)
	c := a + b
	return []interface{}{c}
}
