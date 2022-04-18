package interpreter

import "fmt"

// 用于 native function 单元测试
func printChar(args []interface{}) []interface{} {
	fmt.Printf("%c", args[0].(int32))
	return nil
}

// 用于 native function 单元测试
func printInt(args []interface{}) []interface{} {
	fmt.Printf("%d", args[0].(int32))
	return nil
}

// 用于 native function 单元测试
func add_i32(args []interface{}) []interface{} {
	a := args[0].(int32)
	b := args[1].(int32)
	c := a + b
	return []interface{}{c}
}
