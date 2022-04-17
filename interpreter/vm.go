package interpreter

import (
	"wasmvm/binary"
)

type vm struct {
	operandStack operandStack
	controlStack controlStack
	module       binary.Module
	memory       *memory

	// 全局变量表
	globals []*globalVar

	// 记录第一个局部变量（包括函数参数）在栈中的位置，用于
	// 方便按索引访问栈中的局部变量，它的值等于从栈顶开始
	// 第一个 `函数调用帧` 的 BP（base pointer）
	// 对于函数内的流程控制所产生帧，不更新 local0Idx 的值。
	local0Idx uint32 // name: currentCallFrameBasePointer

	// 注：
	// 目前局部变量（包括函数参数）表直接在操作栈中实现
}

func (v *vm) enterBlock(opcode byte, func_type binary.FuncType,
	instructions []binary.Instruction) {
	bp := v.operandStack.stackSize() - len(func_type.ParamTypes)
	frame := newControlFrame(opcode, func_type, instructions, bp)
	v.controlStack.pushControlFrame(frame)

	if opcode == binary.Call {
		v.local0Idx = uint32(bp)
	}
}

func (v *vm) exitBlock() { // name: leaveBlock
	frame := v.controlStack.popControlFrame() // 消掉当前控制帧
	v.clearBlock(frame)                       // 做一些离开 `被调用者` 之后的清理工作
}

// todo:: 考虑把 clearBlock 函数合并到 exitBlock
func (v *vm) clearBlock(frame *controlFrame) {
	// 这里的 controlFrame 是退出的 `目标层`，而不是 `源层`

	// 丢弃自当前函数 bp (base pointer) 以后产生的所有操作数槽，防止 `被调用者` 产生的
	// 残留数据。
	residues := v.operandStack.stackSize() - frame.bp - len(frame.bt.ResultTypes)

	if residues > 0 {
		// 先弹出有用的数据（即返回值）
		returnValues := v.operandStack.popValues(len(frame.bt.ResultTypes))
		// 丢弃残留数据
		v.operandStack.popValues(residues)
		// 再压入有用的数据
		v.operandStack.pushValues(returnValues)
	}

	// 如果当前是函数退出，则还需要
	// 更新 local0Idx 的值
	if frame.opcode == binary.Call &&
		v.controlStack.controlDepth() > 0 {
		lastCallFrame := v.controlStack.topCallFrame()
		v.local0Idx = uint32(lastCallFrame.bp)
	}

}

func (v *vm) resetBlock(frame *controlFrame) {
	// 注意这里要弹出目标层参数所需数量的操作数，而不是 `源层` 的返回值数量的操作数。
	targetBlockArguments := v.operandStack.popValues(len(frame.bt.ParamTypes))

	// 丢弃目标层 bp 到栈顶的数据
	v.operandStack.popValues(v.operandStack.stackSize() - frame.bp)

	v.operandStack.pushValues(targetBlockArguments)
}

type instructionExecFunc = func(v *vm, args interface{})

// 指令的解析/执行函数表
var instructionTable = make([]instructionExecFunc, 256)

func (v *vm) initMem() {
	// 当前 wasm 只支持创建一个内存块
	if len(v.module.MemSec) != 0 {
		v.memory = newMemory(v.module.MemSec[0])
	}

	// 读取 Data 段，初始化内存块的内容
	for _, dataItem := range v.module.DataSec {
		// 执行偏移值表达式（通常是一个 i32.const 指令）
		for _, offsetInst := range dataItem.Offset {
			v.execInstruction(offsetInst)
		}

		// 操作数栈的顶端操作数————即偏移值表达式的运算结果————表示内存的有效地址
		eaddr := v.operandStack.popU64() // 有效地址是 33 位的无符号整数，这里使用 uint64 来存储
		v.memory.Write(eaddr, dataItem.Init)
	}
}

func (v *vm) initMemWithInitData(init_data []byte) {
	v.memory = newMemoryWithInitData(init_data)
}

func ExecMainFunc(module binary.Module) {
	func_idx := uint32(*module.StartSec) - uint32(len(module.ImportSec)) // 导入函数也会占用函数索引
	v := &vm{module: module}
	v.execFunc(func_idx)
}

// 执行指定函数
// 返回操作数栈和内容的内容，用于测试
func TestFunc(module binary.Module, func_idx uint32) ([]uint64, []byte) {
	v := &vm{module: module}
	v.initMem()
	v.execFunc(func_idx)
	return v.operandStack.slots, dumpMemory(v.memory)
}

func TestFuncWithInitMemoryData(module binary.Module, init_data []byte, func_idx uint32) ([]uint64, []byte) {
	v := &vm{module: module}
	v.initMemWithInitData(init_data)
	v.execFunc(func_idx)
	return v.operandStack.slots, dumpMemory(v.memory)
}

func dumpMemory(m *memory) []byte {
	if m == nil {
		return []byte{}
	} else {
		return m.data
	}
}

// 执行指定函数
func (v *vm) execFunc(func_idx uint32) {
	//code := v.module.CodeSec[func_idx]
	// for _, inst := range code.Expr {
	// 	v.execInst(inst)
	// }

	call(v, func_idx)

	// 	v.loop()
	// }
	// func (v *vm) loop() {

	// 程序的入口是一个模块内部的用户自定义函数，调用 call 方法之后，控制栈
	// 应该有 1 个栈帧，所以这里的 depth 值为 1
	// todo::
	// 按理说下面的循环可以简化成
	// `for v.controlStack.controlDepth()  >= 1 {`

	depth := v.controlStack.controlDepth()
	for v.controlStack.controlDepth() >= depth {
		frame := v.controlStack.topControlFrame()
		if frame.pc == len(frame.instructions) {
			v.exitBlock()
		} else {
			instr := frame.instructions[frame.pc]
			frame.pc++ // 向前移动一个指令
			v.execInstruction(instr)
		}
	}
}

// 执行一条指令
func (v *vm) execInstruction(inst binary.Instruction) {
	instructionTable[inst.Opcode](v, inst.Args)
}

func init() {
	// 控制指令
	instructionTable[binary.Unreachable] = unreachable
	instructionTable[binary.Nop] = nop
	instructionTable[binary.Block] = block
	instructionTable[binary.Loop] = loop
	instructionTable[binary.If] = if_
	instructionTable[binary.Br] = br
	instructionTable[binary.BrIf] = brIf
	instructionTable[binary.BrTable] = brTable
	instructionTable[binary.Return] = return_

	// 操作数（参数 parameter）指令
	instructionTable[binary.Drop] = drop
	instructionTable[binary.Select] = select_

	// 数值指令
	instructionTable[binary.I32Const] = i32Const
	instructionTable[binary.I64Const] = i64Const
	instructionTable[binary.F32Const] = f32Const
	instructionTable[binary.F64Const] = f64Const
	instructionTable[binary.I32Eqz] = i32Eqz
	instructionTable[binary.I32Eq] = i32Eq
	instructionTable[binary.I32Ne] = i32Ne
	instructionTable[binary.I32LtS] = i32LtS
	instructionTable[binary.I32LtU] = i32LtU
	instructionTable[binary.I32GtS] = i32GtS
	instructionTable[binary.I32GtU] = i32GtU
	instructionTable[binary.I32LeS] = i32LeS
	instructionTable[binary.I32LeU] = i32LeU
	instructionTable[binary.I32GeS] = i32GeS
	instructionTable[binary.I32GeU] = i32GeU
	instructionTable[binary.I64Eqz] = i64Eqz
	instructionTable[binary.I64Eq] = i64Eq
	instructionTable[binary.I64Ne] = i64Ne
	instructionTable[binary.I64LtS] = i64LtS
	instructionTable[binary.I64LtU] = i64LtU
	instructionTable[binary.I64GtS] = i64GtS
	instructionTable[binary.I64GtU] = i64GtU
	instructionTable[binary.I64LeS] = i64LeS
	instructionTable[binary.I64LeU] = i64LeU
	instructionTable[binary.I64GeS] = i64GeS
	instructionTable[binary.I64GeU] = i64GeU
	instructionTable[binary.F32Eq] = f32Eq
	instructionTable[binary.F32Ne] = f32Ne
	instructionTable[binary.F32Lt] = f32Lt
	instructionTable[binary.F32Gt] = f32Gt
	instructionTable[binary.F32Le] = f32Le
	instructionTable[binary.F32Ge] = f32Ge
	instructionTable[binary.F64Eq] = f64Eq
	instructionTable[binary.F64Ne] = f64Ne
	instructionTable[binary.F64Lt] = f64Lt
	instructionTable[binary.F64Gt] = f64Gt
	instructionTable[binary.F64Le] = f64Le
	instructionTable[binary.F64Ge] = f64Ge
	instructionTable[binary.I32Clz] = i32Clz
	instructionTable[binary.I32Ctz] = i32Ctz
	instructionTable[binary.I32PopCnt] = i32PopCnt
	instructionTable[binary.I32Add] = i32Add
	instructionTable[binary.I32Sub] = i32Sub
	instructionTable[binary.I32Mul] = i32Mul
	instructionTable[binary.I32DivS] = i32DivS
	instructionTable[binary.I32DivU] = i32DivU
	instructionTable[binary.I32RemS] = i32RemS
	instructionTable[binary.I32RemU] = i32RemU
	instructionTable[binary.I32And] = i32And
	instructionTable[binary.I32Or] = i32Or
	instructionTable[binary.I32Xor] = i32Xor
	instructionTable[binary.I32Shl] = i32Shl
	instructionTable[binary.I32ShrS] = i32ShrS
	instructionTable[binary.I32ShrU] = i32ShrU
	instructionTable[binary.I32Rotl] = i32Rotl
	instructionTable[binary.I32Rotr] = i32Rotr
	instructionTable[binary.I64Clz] = i64Clz
	instructionTable[binary.I64Ctz] = i64Ctz
	instructionTable[binary.I64PopCnt] = i64PopCnt
	instructionTable[binary.I64Add] = i64Add
	instructionTable[binary.I64Sub] = i64Sub
	instructionTable[binary.I64Mul] = i64Mul
	instructionTable[binary.I64DivS] = i64DivS
	instructionTable[binary.I64DivU] = i64DivU
	instructionTable[binary.I64RemS] = i64RemS
	instructionTable[binary.I64RemU] = i64RemU
	instructionTable[binary.I64And] = i64And
	instructionTable[binary.I64Or] = i64Or
	instructionTable[binary.I64Xor] = i64Xor
	instructionTable[binary.I64Shl] = i64Shl
	instructionTable[binary.I64ShrS] = i64ShrS
	instructionTable[binary.I64ShrU] = i64ShrU
	instructionTable[binary.I64Rotl] = i64Rotl
	instructionTable[binary.I64Rotr] = i64Rotr
	instructionTable[binary.F32Abs] = f32Abs
	instructionTable[binary.F32Neg] = f32Neg
	instructionTable[binary.F32Ceil] = f32Ceil
	instructionTable[binary.F32Floor] = f32Floor
	instructionTable[binary.F32Trunc] = f32Trunc
	instructionTable[binary.F32Nearest] = f32Nearest
	instructionTable[binary.F32Sqrt] = f32Sqrt
	instructionTable[binary.F32Add] = f32Add
	instructionTable[binary.F32Sub] = f32Sub
	instructionTable[binary.F32Mul] = f32Mul
	instructionTable[binary.F32Div] = f32Div
	instructionTable[binary.F32Min] = f32Min
	instructionTable[binary.F32Max] = f32Max
	instructionTable[binary.F32CopySign] = f32CopySign
	instructionTable[binary.F64Abs] = f64Abs
	instructionTable[binary.F64Neg] = f64Neg
	instructionTable[binary.F64Ceil] = f64Ceil
	instructionTable[binary.F64Floor] = f64Floor
	instructionTable[binary.F64Trunc] = f64Trunc
	instructionTable[binary.F64Nearest] = f64Nearest
	instructionTable[binary.F64Sqrt] = f64Sqrt
	instructionTable[binary.F64Add] = f64Add
	instructionTable[binary.F64Sub] = f64Sub
	instructionTable[binary.F64Mul] = f64Mul
	instructionTable[binary.F64Div] = f64Div
	instructionTable[binary.F64Min] = f64Min
	instructionTable[binary.F64Max] = f64Max
	instructionTable[binary.F64CopySign] = f64CopySign
	instructionTable[binary.I32WrapI64] = i32WrapI64
	instructionTable[binary.I32TruncF32S] = i32TruncF32S
	instructionTable[binary.I32TruncF32U] = i32TruncF32U
	instructionTable[binary.I32TruncF64S] = i32TruncF64S
	instructionTable[binary.I32TruncF64U] = i32TruncF64U
	instructionTable[binary.I64ExtendI32S] = i64ExtendI32S
	instructionTable[binary.I64ExtendI32U] = i64ExtendI32U
	instructionTable[binary.I64TruncF32S] = i64TruncF32S
	instructionTable[binary.I64TruncF32U] = i64TruncF32U
	instructionTable[binary.I64TruncF64S] = i64TruncF64S
	instructionTable[binary.I64TruncF64U] = i64TruncF64U
	instructionTable[binary.F32ConvertI32S] = f32ConvertI32S
	instructionTable[binary.F32ConvertI32U] = f32ConvertI32U
	instructionTable[binary.F32ConvertI64S] = f32ConvertI64S
	instructionTable[binary.F32ConvertI64U] = f32ConvertI64U
	instructionTable[binary.F32DemoteF64] = f32DemoteF64
	instructionTable[binary.F64ConvertI32S] = f64ConvertI32S
	instructionTable[binary.F64ConvertI32U] = f64ConvertI32U
	instructionTable[binary.F64ConvertI64S] = f64ConvertI64S
	instructionTable[binary.F64ConvertI64U] = f64ConvertI64U
	instructionTable[binary.F64PromoteF32] = f64PromoteF32
	instructionTable[binary.I32ReinterpretF32] = i32ReinterpretF32
	instructionTable[binary.I64ReinterpretF64] = i64ReinterpretF64
	instructionTable[binary.F32ReinterpretI32] = f32ReinterpretI32
	instructionTable[binary.F64ReinterpretI64] = f64ReinterpretI64
	instructionTable[binary.I32Extend8S] = i32Extend8S
	instructionTable[binary.I32Extend16S] = i32Extend16S
	instructionTable[binary.I64Extend8S] = i64Extend8S
	instructionTable[binary.I64Extend16S] = i64Extend16S
	instructionTable[binary.I64Extend32S] = i64Extend32S
	instructionTable[binary.TruncSat] = truncSat

	// 内存指令
	instructionTable[binary.I32Load] = i32Load
	instructionTable[binary.I64Load] = i64Load
	instructionTable[binary.F32Load] = f32Load
	instructionTable[binary.F64Load] = f64Load
	instructionTable[binary.I32Load8S] = i32Load8S
	instructionTable[binary.I32Load8U] = i32Load8U
	instructionTable[binary.I32Load16S] = i32Load16S
	instructionTable[binary.I32Load16U] = i32Load16U
	instructionTable[binary.I64Load8S] = i64Load8S
	instructionTable[binary.I64Load8U] = i64Load8U
	instructionTable[binary.I64Load16S] = i64Load16S
	instructionTable[binary.I64Load16U] = i64Load16U
	instructionTable[binary.I64Load32S] = i64Load32S
	instructionTable[binary.I64Load32U] = i64Load32U
	instructionTable[binary.I32Store] = i32Store
	instructionTable[binary.I64Store] = i64Store
	instructionTable[binary.F32Store] = f32Store
	instructionTable[binary.F64Store] = f64Store
	instructionTable[binary.I32Store8] = i32Store8
	instructionTable[binary.I32Store16] = i32Store16
	instructionTable[binary.I64Store8] = i64Store8
	instructionTable[binary.I64Store16] = i64Store16
	instructionTable[binary.I64Store32] = i64Store32
	instructionTable[binary.MemorySize] = memorySize
	instructionTable[binary.MemoryGrow] = memoryGrow

	// 函数指令
	instructionTable[binary.Call] = call

	// 变量指令
	instructionTable[binary.LocalGet] = localGet
	instructionTable[binary.LocalSet] = localSet
	instructionTable[binary.LocalTee] = localTee
	instructionTable[binary.GlobalGet] = globalGet
	instructionTable[binary.GlobalSet] = globalSet
}
