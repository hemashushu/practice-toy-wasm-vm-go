package interpreter

import "wasmvm/binary"

type vm struct {
	operandStack operandStack
	module       binary.Module
}

type instFn = func(v *vm, args interface{})

// 指令的解析/执行函数表
var instTable = make([]instFn, 256)

func init() {
	instTable[binary.Call] = call // hack!
	instTable[binary.Drop] = drop
	instTable[binary.Select] = select_
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
}

func ExecMainFunc(module binary.Module) {
	func_idx := uint32(*module.StartSec) - uint32(len(module.ImportSec)) // 导入函数也会占用函数索引
	v := &vm{module: module}
	v.execCode(func_idx)
}

// 执行指定函数
// 返回操作数栈的内容，用于测试
func ExecFunc(module binary.Module, func_idx uint32) []uint64 {
	v := &vm{module: module}
	v.execCode(func_idx)
	return v.operandStack.slots
}

// 执行指定函数
func (v *vm) execCode(func_idx uint32) {
	code := v.module.CodeSec[func_idx]
	for _, inst := range code.Expr {
		v.execInst(inst)
	}
}

// 执行一条指令
func (v *vm) execInst(inst binary.Instruction) {
	instTable[inst.Opcode](v, inst.Args)
}
