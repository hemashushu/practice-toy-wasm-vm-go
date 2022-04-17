package interpreter

import (
	"errors"
)

// ======== 控制指令

// unreachable 标记当前位置不应该到达，如果这个指令被执行，抛出不可恢复异常
// nop 空指令
//
// https://developer.mozilla.org/en-US/docs/WebAssembly/Reference/Control_flow/unreachable
// https://developer.mozilla.org/en-US/docs/WebAssembly/Reference/Control_flow/nop

func unreachable(v *vm, _ interface{}) {
	panic(errors.New("unreachable"))
}

func nop(v *vm, _ interface{}) {
	//
}
