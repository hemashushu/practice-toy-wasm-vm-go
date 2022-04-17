package interpreter

import "wasmvm/binary"

// ======== 控制指令

// -------- 流程控制指令
// block
// loop
// if
// else
// end

func block(v *vm, args interface{}) {
	// 可以对比 callInteralFunc() 函数的实现，基本上一样

	blockArgs := args.(binary.BlockArgs)

	// 比较合适的名称：
	// funcType := v.module.convertBlockTypeIntoFunctionType(blockArgs.BT)
	funcType := v.module.GetBlockType(blockArgs.BT)
	expr := blockArgs.Instrs

	v.enterBlock(binary.Block, funcType, expr)

	// 注：block 没有自己的局部变量空槽
}

func loop(v *vm, args interface{}) {
	// loop() 函数基本上跟 block() 一样

	blockArgs := args.(binary.BlockArgs)

	funcType := v.module.GetBlockType(blockArgs.BT)
	expr := blockArgs.Instrs

	v.enterBlock(binary.Loop, funcType, expr)
}

func if_(v *vm, args interface{}) {
	ifArgs := args.(binary.IfArgs)

	funcType := v.module.GetBlockType(ifArgs.BT)

	// if 结构的两个分支共用同一个 block type
	if v.operandStack.popBool() {
		v.enterBlock(binary.If, funcType, ifArgs.Instrs1)
	} else {
		v.enterBlock(binary.If, funcType, ifArgs.Instrs2)
	}
}

// -------- 分支/跳转指令
// br
// br_if
// br_table
// return

// 函数调用
// call
// call_indirect

// -------- br 指令
//
// br 指令后面接着 `跳转目标` 的相对深度。
// 对于 block/if 指令来说，跳转目标是指令的结尾处（即 end 指令），
// 对于 loop 指令来说，跳转目标是指令的开始处（即 loop 指令）。
//
// (func
// 	(block
// 		(i32.const 100)
// 		(br 0)		        ;; 跳转目标为 dest_a
// 		(i32.const 101)
// 	)						;; dest_a
// 	(loop                   ;; dest_b
// 		(i32.const 200)
// 		(br 0)				;; 跳转目标为 dest_b
// 		(i32.const 201)
// 	)
// 	(if (i32.eqz (i32.const 300))
// 		(then (i32.const 400) (br 0) (i32.const 401))	;; 跳转目标为 dest_b
// 		(else (i32.const 500) (br 0) (i32.const 501))	;; 跳转目标为 dest_b
// 	)						;; dest_c
// )
//
// "br 指令对于 `控制块`" 跟
// "return 指令对于 `函数`" 的处理方式是一样

func br(v *vm, args interface{}) {
	relative_depth := int(args.(uint32))

	// 先放心弹出 relative_depth 个控制帧
	// 注：
	// 该数值从 0 开始，当 relative_depth 为 N 时，实际上有 N+1 层
	// 结构块，先直接弹出 N 层，最后一层则需判断是 block/if 还是 loop
	// 指令，前者则跳到结构末尾，后者则跳到结构开头。

	for i := 0; i < relative_depth; i++ {
		v.controlStack.popControlFrame()
	}

	if frame := v.controlStack.topControlFrame(); frame.opcode == binary.Loop {
		// 目标层是一个 loop 结构
		// 需要弹出目标层参数所需数量的操作数，然后跳到目标层的第一行指令
		v.resetBlock(frame)
		frame.pc = 0 // todo:: 考虑把这个方法挪进 resetBlock 函数
	} else {
		// 目标层是一个 if 或者 block 结构，或者是函数
		v.exitBlock()
	}
}

func brIf(v *vm, args interface{}) {
	// br_if 指令先从操作数栈顶弹出一个有符号的整数（int32），非 0 则执行 br 操作，
	// 等于 0 则什么都不做（仅仅消耗掉栈顶的一个操作数）
	if v.operandStack.popBool() {
		br(v, args)
	}
}

func brTable(v *vm, args interface{}) {
	// br_table 指令先从操作数栈顶弹出一个 uint32 整数，这个数将作为
	// br_table 后面的整数列表的索引，获取跳转的目标。如果该索引超出了
	// 列表范围，则跳转目标的 br_table 指令的最末尾一个参数（即默认目标）

	brTableArgs := args.(binary.BrTableArgs)
	idx := int(v.operandStack.popU32())
	if idx < len(brTableArgs.Labels) {
		br(v, brTableArgs.Labels[idx])
	} else {
		br(v, brTableArgs.Default)
	}
}

func return_(v *vm, _ interface{}) {
	blockRelativeDepth := v.controlStack.getRelativeDepth()
	br(v, uint32(blockRelativeDepth))
}
