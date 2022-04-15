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

func (v *vm) enterBlock(opcode byte, bt binary.FuncType,
	instructions []binary.Instruction) {
	bp := v.operandStack.stackSize() - len(bt.ParamTypes)
	frame := newControlFrame(opcode, bt, instructions, bp)
	v.controlStack.pushControlFrame(frame)

	if opcode == binary.Call {
		v.local0Idx = uint32(bp)
	}
}

func (v *vm) exitBlock() { // name: leaveBlock
	frame := v.controlStack.popControlFrame() // 消掉当前控制帧
	v.clearBlock(frame)                       // 做一些离开 `被调用者` 之后的清理工作
}

func (v *vm) clearBlock(frame *controlFrame) {
	// 丢弃自当前函数 bp (base pointer) 以后产生的所有操作数槽，防止 `被调用者` 产生的
	// 残留数据。
	residues := v.operandStack.stackSize() - frame.bp - len(frame.bt.ResultTypes)

	if residues > 0 {
		// 先弹出有用的数据（即返回值）
		returnValues := v.operandStack.popU64Values(len(frame.bt.ResultTypes))
		// 丢弃残留数据
		v.operandStack.popU64Values(residues)
		// 再压入有用的数据
		v.operandStack.pushU64Values(returnValues)
	}

	// 如果当前是函数退出，则还需要
	// 更新 local0Idx 的值
	if frame.opcode == binary.Call &&
		v.controlStack.controlDepth() > 0 {
		lastCallFrame, _ := v.controlStack.topCallFrame()
		v.local0Idx = uint32(lastCallFrame.bp)
	}

}

type instFn = func(v *vm, args interface{})

// 指令的解析/执行函数表
var instTable = make([]instFn, 256)

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
	instTable[inst.Opcode](v, inst.Args)
}

func init() {
	// 操作数（参数 parameter）指令
	instTable[binary.Drop] = drop
	instTable[binary.Select] = select_

	// 数值指令
	instTable[binary.I32Const] = i32Const
	instTable[binary.I64Const] = i64Const
	instTable[binary.F32Const] = f32Const
	instTable[binary.F64Const] = f64Const
	instTable[binary.I32Eqz] = i32Eqz
	instTable[binary.I32Eq] = i32Eq
	instTable[binary.I32Ne] = i32Ne
	instTable[binary.I32LtS] = i32LtS
	instTable[binary.I32LtU] = i32LtU
	instTable[binary.I32GtS] = i32GtS
	instTable[binary.I32GtU] = i32GtU
	instTable[binary.I32LeS] = i32LeS
	instTable[binary.I32LeU] = i32LeU
	instTable[binary.I32GeS] = i32GeS
	instTable[binary.I32GeU] = i32GeU
	instTable[binary.I64Eqz] = i64Eqz
	instTable[binary.I64Eq] = i64Eq
	instTable[binary.I64Ne] = i64Ne
	instTable[binary.I64LtS] = i64LtS
	instTable[binary.I64LtU] = i64LtU
	instTable[binary.I64GtS] = i64GtS
	instTable[binary.I64GtU] = i64GtU
	instTable[binary.I64LeS] = i64LeS
	instTable[binary.I64LeU] = i64LeU
	instTable[binary.I64GeS] = i64GeS
	instTable[binary.I64GeU] = i64GeU
	instTable[binary.F32Eq] = f32Eq
	instTable[binary.F32Ne] = f32Ne
	instTable[binary.F32Lt] = f32Lt
	instTable[binary.F32Gt] = f32Gt
	instTable[binary.F32Le] = f32Le
	instTable[binary.F32Ge] = f32Ge
	instTable[binary.F64Eq] = f64Eq
	instTable[binary.F64Ne] = f64Ne
	instTable[binary.F64Lt] = f64Lt
	instTable[binary.F64Gt] = f64Gt
	instTable[binary.F64Le] = f64Le
	instTable[binary.F64Ge] = f64Ge
	instTable[binary.I32Clz] = i32Clz
	instTable[binary.I32Ctz] = i32Ctz
	instTable[binary.I32PopCnt] = i32PopCnt
	instTable[binary.I32Add] = i32Add
	instTable[binary.I32Sub] = i32Sub
	instTable[binary.I32Mul] = i32Mul
	instTable[binary.I32DivS] = i32DivS
	instTable[binary.I32DivU] = i32DivU
	instTable[binary.I32RemS] = i32RemS
	instTable[binary.I32RemU] = i32RemU
	instTable[binary.I32And] = i32And
	instTable[binary.I32Or] = i32Or
	instTable[binary.I32Xor] = i32Xor
	instTable[binary.I32Shl] = i32Shl
	instTable[binary.I32ShrS] = i32ShrS
	instTable[binary.I32ShrU] = i32ShrU
	instTable[binary.I32Rotl] = i32Rotl
	instTable[binary.I32Rotr] = i32Rotr
	instTable[binary.I64Clz] = i64Clz
	instTable[binary.I64Ctz] = i64Ctz
	instTable[binary.I64PopCnt] = i64PopCnt
	instTable[binary.I64Add] = i64Add
	instTable[binary.I64Sub] = i64Sub
	instTable[binary.I64Mul] = i64Mul
	instTable[binary.I64DivS] = i64DivS
	instTable[binary.I64DivU] = i64DivU
	instTable[binary.I64RemS] = i64RemS
	instTable[binary.I64RemU] = i64RemU
	instTable[binary.I64And] = i64And
	instTable[binary.I64Or] = i64Or
	instTable[binary.I64Xor] = i64Xor
	instTable[binary.I64Shl] = i64Shl
	instTable[binary.I64ShrS] = i64ShrS
	instTable[binary.I64ShrU] = i64ShrU
	instTable[binary.I64Rotl] = i64Rotl
	instTable[binary.I64Rotr] = i64Rotr
	instTable[binary.F32Abs] = f32Abs
	instTable[binary.F32Neg] = f32Neg
	instTable[binary.F32Ceil] = f32Ceil
	instTable[binary.F32Floor] = f32Floor
	instTable[binary.F32Trunc] = f32Trunc
	instTable[binary.F32Nearest] = f32Nearest
	instTable[binary.F32Sqrt] = f32Sqrt
	instTable[binary.F32Add] = f32Add
	instTable[binary.F32Sub] = f32Sub
	instTable[binary.F32Mul] = f32Mul
	instTable[binary.F32Div] = f32Div
	instTable[binary.F32Min] = f32Min
	instTable[binary.F32Max] = f32Max
	instTable[binary.F32CopySign] = f32CopySign
	instTable[binary.F64Abs] = f64Abs
	instTable[binary.F64Neg] = f64Neg
	instTable[binary.F64Ceil] = f64Ceil
	instTable[binary.F64Floor] = f64Floor
	instTable[binary.F64Trunc] = f64Trunc
	instTable[binary.F64Nearest] = f64Nearest
	instTable[binary.F64Sqrt] = f64Sqrt
	instTable[binary.F64Add] = f64Add
	instTable[binary.F64Sub] = f64Sub
	instTable[binary.F64Mul] = f64Mul
	instTable[binary.F64Div] = f64Div
	instTable[binary.F64Min] = f64Min
	instTable[binary.F64Max] = f64Max
	instTable[binary.F64CopySign] = f64CopySign
	instTable[binary.I32WrapI64] = i32WrapI64
	instTable[binary.I32TruncF32S] = i32TruncF32S
	instTable[binary.I32TruncF32U] = i32TruncF32U
	instTable[binary.I32TruncF64S] = i32TruncF64S
	instTable[binary.I32TruncF64U] = i32TruncF64U
	instTable[binary.I64ExtendI32S] = i64ExtendI32S
	instTable[binary.I64ExtendI32U] = i64ExtendI32U
	instTable[binary.I64TruncF32S] = i64TruncF32S
	instTable[binary.I64TruncF32U] = i64TruncF32U
	instTable[binary.I64TruncF64S] = i64TruncF64S
	instTable[binary.I64TruncF64U] = i64TruncF64U
	instTable[binary.F32ConvertI32S] = f32ConvertI32S
	instTable[binary.F32ConvertI32U] = f32ConvertI32U
	instTable[binary.F32ConvertI64S] = f32ConvertI64S
	instTable[binary.F32ConvertI64U] = f32ConvertI64U
	instTable[binary.F32DemoteF64] = f32DemoteF64
	instTable[binary.F64ConvertI32S] = f64ConvertI32S
	instTable[binary.F64ConvertI32U] = f64ConvertI32U
	instTable[binary.F64ConvertI64S] = f64ConvertI64S
	instTable[binary.F64ConvertI64U] = f64ConvertI64U
	instTable[binary.F64PromoteF32] = f64PromoteF32
	instTable[binary.I32ReinterpretF32] = i32ReinterpretF32
	instTable[binary.I64ReinterpretF64] = i64ReinterpretF64
	instTable[binary.F32ReinterpretI32] = f32ReinterpretI32
	instTable[binary.F64ReinterpretI64] = f64ReinterpretI64
	instTable[binary.I32Extend8S] = i32Extend8S
	instTable[binary.I32Extend16S] = i32Extend16S
	instTable[binary.I64Extend8S] = i64Extend8S
	instTable[binary.I64Extend16S] = i64Extend16S
	instTable[binary.I64Extend32S] = i64Extend32S
	instTable[binary.TruncSat] = truncSat

	// 内存指令
	instTable[binary.I32Load] = i32Load
	instTable[binary.I64Load] = i64Load
	instTable[binary.F32Load] = f32Load
	instTable[binary.F64Load] = f64Load
	instTable[binary.I32Load8S] = i32Load8S
	instTable[binary.I32Load8U] = i32Load8U
	instTable[binary.I32Load16S] = i32Load16S
	instTable[binary.I32Load16U] = i32Load16U
	instTable[binary.I64Load8S] = i64Load8S
	instTable[binary.I64Load8U] = i64Load8U
	instTable[binary.I64Load16S] = i64Load16S
	instTable[binary.I64Load16U] = i64Load16U
	instTable[binary.I64Load32S] = i64Load32S
	instTable[binary.I64Load32U] = i64Load32U
	instTable[binary.I32Store] = i32Store
	instTable[binary.I64Store] = i64Store
	instTable[binary.F32Store] = f32Store
	instTable[binary.F64Store] = f64Store
	instTable[binary.I32Store8] = i32Store8
	instTable[binary.I32Store16] = i32Store16
	instTable[binary.I64Store8] = i64Store8
	instTable[binary.I64Store16] = i64Store16
	instTable[binary.I64Store32] = i64Store32
	instTable[binary.MemorySize] = memorySize
	instTable[binary.MemoryGrow] = memoryGrow

	// 函数指令
	instTable[binary.Call] = call // hack!

}
