package binary

// 值类型
//
// 虚拟机只支持 4 种类型：i32, i64, f32, f64
//
// 其中 i32， i64 在内部是指 `无符号的整数`，
// 但部分指令是需要将 `无符号整数` 解析为 `有符号整数` 再进行运算，
// 比如 `lt_u`, `lt_s`，所有有些地方会出现 u32, u64, s32, s64 等名称。
//
// https://webassembly.github.io/spec/core/syntax/values.html#integers

type ValType = byte

const (
	ValTypeI32 ValType = 0x7f // i32
	ValTypeI64 ValType = 0x7E // i64
	ValTypeF32 ValType = 0x7D // f32
	ValTypeF64 ValType = 0x7C // f64
)

// 指令块（的返回值）类型

type BlockType = int32 // leb128 编码

const (
	BlockTypeI32   BlockType = -1  // ()->(i32)
	BlockTypeI64   BlockType = -2  // ()->(i64)
	BlockTypeF32   BlockType = -3  // ()->(f32)
	BlockTypeF64   BlockType = -4  // ()->(f64)
	BlockTypeEmpty BlockType = -64 // ()->()
)

type Expr = []Instruction

type Instruction struct {
	Opcode byte
	Args   interface{}
}

// ================ 指令详细
//
// 指令（字节码）存在的地方:
// 1. 全局项的初始值表达式
// 2. 元素项的偏移值表达式（元素项用于构建表的初始内容）
// 3. 数据项的偏移值表达式（数据项用于构建内存块的初始内容）
// 4. 代码项的字节码
//
// global: global_type + init_expr
// elem:   table_idx + offset_expr + <func_idx>
// data:   mem_block_idx + offset_expr + <byte>
// code:   byte_count + <locals> + expr
//
// expr: inst* + 0x0b

// ---------------- 数值指令
//
// i32.const: 0x41 + i32  // 参数是一个 leb128 int32（**有符号**）
// i64.const: 0x42 + i64  // 参数是一个 leb128 int64（**有符号**）
// f32.const: 0x43 + f32  // 参数是一个定长 4 字节 float32
// f64.const: 0x44 + f64  // 参数是一个定长 8 字节 float64
// trunc_sat: 0xFC + byte // 参数是一个 byte
//
// (module
// 	(func
// 		(f32.const 12.3)
// 		(f32.const 45.6)
// 		(f32.add)
// 		(i32.trunc_sat_f32_s)
// 		(drop)
// 	)
// )
//
// 0x0015 | 10          | size of function
// 0x0016 | 00          | 0 local blocks
// 0x0017 | 43 cd cc 44 | F32Const { value: Ieee32(1095027917) }
//        | 41
// 0x001c | 43 66 66 36 | F32Const { value: Ieee32(1110861414) }
//        | 42
// 0x0021 | 92          | F32Add
// 0x0022 | fc 00       | I32TruncSatF32S
// 0x0024 | 1a          | Drop
// 0x0025 | 0b          | End
//
// 函数的指令（字节码）以 0x0B 结束

// ---------------- 变量指令
//
// local.get:   0x20 + local_idx  // 参数是一个 leb128 uint32（无符号），下同
// local.set:   0x21 + local_idx
// local.tee:   0x22 + local_idx
// global.get:  0x23 + global_idx
// global.set : 0x24 + global_idx
//
//
// (module
// 	(global $g1 (mut i32) (i32.const 1))  ;; $g1, $g2 可视为自动索引值
// 	(global $g2 (mut i32) (i32.const 2))
// 	(func (param $a i32) (param $b i32)
// 		(global.get $g1)
// 		(global.set $g2)
// 		(local.get $a)
// 		(local.set $b)
// 	)
// )
//
// 0x0024 | 0a          | size of function
// 0x0025 | 00          | 0 local blocks
// 0x0026 | 23 00       | GlobalGet { global_index: 0 }
// 0x0028 | 24 01       | GlobalSet { global_index: 1 }
// 0x002a | 20 00       | LocalGet { local_index: 0 }
// 0x002c | 21 01       | LocalSet { local_index: 1 }
// 0x002e | 0b          | End

// ---------------- 内存指令
//
// xxx.load?_?:  0x28..0x35 + align:uint32 + offset:uint32
// xxx.store?:   0x36..0x3e + align:uint32 + offset:uint32
// memory.size:  0x3F + mem_block_idx:uint32
// memory.grow:  0x40 + mem_block_idx:uint32
//
// (module
// 	(memory 1 8)
// 	(data (offset (i32.const 100)) "hello")
//
// 	(func
// 		(i32.const 1)
// 		(i32.const 2)
// 		(i32.load offset=100)
// 		(i32.store offset=100)
// 		(memory.size)
// 		(drop)
// 		(i32.const 4)
// 		(memory.grow)
// 		(drop)
// 	)
// )
//
// 0x001b | 14          | size of function
// 0x001c | 00          | 0 local blocks
// 0x001d | 41 01       | I32Const { value: 1 }
// 0x001f | 41 02       | I32Const { value: 2 }
// 0x0021 | 28 02 64    | I32Load { memarg: MemoryImmediate { align: 2, offset: 100, memory: 0 } }
// 0x0024 | 36 02 64    | I32Store { memarg: MemoryImmediate { align: 2, offset: 100, memory: 0 } }
// 0x0027 | 3f 00       | MemorySize { mem: 0, mem_byte: 0 }
// 0x0029 | 1a          | Drop
// 0x002a | 41 04       | I32Const { value: 4 }
// 0x002c | 40 00       | MemoryGrow { mem: 0, mem_byte: 0 }
// 0x002e | 1a          | Drop
// 0x002f | 0b          | End

type MemArg struct {
	Align  uint32
	Offset uint32
}

// ---------------- 结构化控制指令
//
// block:		 0x02 + block_return_type:int32 + inst* + 0x0b
// loop:		 0x03 + block_return_type:int32 + inst* + 0x0b
// if...else...: 0x04 + block_return_type:int32 + inst* + (0x05 + inst*)? + 0x0b
//
// (module
//     (func (result i32)                          ;; (result i32) 是 block return type
//         (block (result i32)                     ;;
//             (i32.const 1)
//             (loop (result i32)                  ;;
//                 (if (result i32) (i32.const 2)  ;; (i32.const 2) 是测试部分
//                     (then (i32.const 3))        ;; then
//                     (else (i32.const 4))        ;; else
//                 )
//             )
//         )
//         (drop)
//     )
// )
//
// 0x0016 | 15          | size of function
// 0x0017 | 00          | 0 local blocks
// 0x0018 | 02 7f       | Block { ty: Type(I32) }
// 0x001a | 41 01       |   I32Const { value: 1 }
// 0x001c | 03 7f       |   Loop { ty: Type(I32) }
// 0x001e | 41 02       |     I32Const { value: 2 }
// 0x0020 | 04 7f       |     If { ty: Type(I32) }
// 0x0022 | 41 03       |       I32Const { value: 3 }
// 0x0024 | 05          |     Else
// 0x0025 | 41 04       |       I32Const { value: 4 }
// 0x0027 | 0b          |     End
// 0x0028 | 0b          |   End
// 0x0029 | 0b          | End
// 0x002a | 1a          | Drop
// 0x002b | 0b          | End

// block 和 loop 指令的结构相同
type BlockArgs struct {
	BT     BlockType // 块的返回值类型
	Instrs []Instruction
}

type IfArgs struct {
	BT      BlockType     // 块的返回值类型
	Instrs1 []Instruction // consequentInstructions
	Instrs2 []Instruction // alternateInstructions
}

// ---------------- 结构块的参数和返回值
//
// 结构化控制指令 block, loop, if 生成的结构也称为 `（结构）块`
// `块` 可视为 "内联函数"，同样允许定义参数和返回值（但不能
// 声明专属于块指令的局部变量）
//
// (func $add (param $a i32) (param $b i32) (result i32)
// 	(local.get $a)
// 	(local.get $b)
// 	(block (param i32 i32) (result i32)
// 		(i32.add)
// 	)
// )
//
// 上面的 block 指令从逻辑上是从操作数栈弹出两个 i32 （因为
// param 指令表明有两个参数），然后创建新的帧栈，然后压入两个操作数。
// 如果具体的实现是用同一个 **操作数栈** 实现帧栈，则其实 block 的
// param 指令什么事都不需要做，仅仅是从逻辑上创建了一个帧栈而已。
//
// loop 指令跟 block 指令，也是可以定义参数和返回值。
//
// if 是 block 的增强版，首先 if 从操作数栈弹出一个数字，
// 非 0 则执行第一个分支（块），0 则执行第二个分支（块），
// 两个块的 "参数和返回值"（即签名）一样
//
// (local.get $operator)
// (if (param i32 i32) (result i32)
//    (then (i32.add))
//    (else (i32.sub))
// )
//

// ---------------- 跳转指令
//
// br: 			0x0c + label_idx						;; 无条件跳转
// br_if: 		0x0d + label_idx						;; 条件跳转
// br_table: 	0x0e + <label_idx> + default_lable_idx  ;; 查表跳转
// return: 		0x0f
//
// return:
// 	 br 的特殊形式，直接跳出到函数首层并返回值
//
// label_idx:
//   所谓 label（`标签`），即跳转的目标，对于 block/if 指令来说，跳转目标是
//   指令的结尾处（即 end 指令），对于 loop 指令来说，跳转目标是
//   指令的开始处（即 loop 指令）。
//
//   标签索引，也就是区块的相对层数
//   吐嘈一下，其实叫做 jump_destination_relative_depth 更合适，
//   或者叫 `相对深度` 更合适
//
//
// (func
// 	(block
// 		(i32.const 100)
// 		(br 0)
// 		(i32.const 101)
// 	)
// 	(loop
// 		(i32.const 200)
// 		(br 0)
// 		(i32.const 201)
// 	)
// 	(if (i32.eqz (i32.const 300))
// 		(then (i32.const 400) (br 0) (i32.const 401))
// 		(else (i32.const 500) (br 0) (i32.const 501))
// 	)
// )
//
// 0x0016 | 30          | size of function
// 0x0017 | 00          | 0 local blocks
// 0x0018 | 02 40       |   Block { ty: Empty }
// 0x001a | 41 e4 00    |     I32Const { value: 100 }
// 0x001d | 0c 00       |     Br { relative_depth: 0 }
// 0x001f | 41 e5 00    |     I32Const { value: 101 }
// 0x0022 | 0b          |   End
// 0x0023 | 03 40       |   Loop { ty: Empty }
// 0x0025 | 41 c8 01    |     I32Const { value: 200 }
// 0x0028 | 0c 00       |     Br { relative_depth: 0 }
// 0x002a | 41 c9 01    |     I32Const { value: 201 }
// 0x002d | 0b          |   End
// 0x002e | 41 ac 02    |   I32Const { value: 300 }
// 0x0031 | 45          |   I32Eqz
// 0x0032 | 04 40       |   If { ty: Empty }
// 0x0034 | 41 90 03    |     I32Const { value: 400 }
// 0x0037 | 0c 00       |     Br { relative_depth: 0 }
// 0x0039 | 41 91 03    |     I32Const { value: 401 }
// 0x003c | 05          |   Else
// 0x003d | 41 f4 03    |     I32Const { value: 500 }
// 0x0040 | 0c 00       |     Br { relative_depth: 0 }
// 0x0042 | 41 f5 03    |     I32Const { value: 501 }
// 0x0045 | 0b          |   End
// 0x0046 | 0b          | End

// -------- 文本格式的 br 指令的跳转标签
//
// br 指令后面除了可以跟一个 uint32 数字，也可以跟一个标识符（`$xx`）
//
// (func
// 	(block $foo									;; 为 block 命名为 $foo
// 		(block $bar								;; 为 block 命名为 $bar
// 			(block $baz							;; 为 block 命名为 $baz
// 				(br 0) (br 1) (br 2)
// 				(br $baz) (br $bar) (br $foo)	;; 跟上一行等同的指令
// 			)
// 		)
// 	)
// )
//
// 使用 return 指令可以直接跳到函数末尾（函数末尾也是一个 `end` 指令）

// -------- br_if br_table 指令
//
// (func
// 	(block
// 		(block
// 			(block
// 				(br 1)
// 				(br_if 2 (i32.const 100))
// 				(br_table 0 1 2 3) ;; 3 是默认标签
// 				(return)
// 			)
// 		)
// 	)
// )
//
// 0x0047 | 19          | size of function
// 0x0048 | 00          | 0 local blocks
// 0x0049 | 02 40       | Block { ty: Empty }
// 0x004b | 02 40       |   Block { ty: Empty }
// 0x004d | 02 40       |     Block { ty: Empty }
// 0x004f | 0c 01       |       Br { relative_depth: 1 }
// 0x0051 | 41 e4 00    |       I32Const { value: 100 }
// 0x0054 | 0d 02       |       BrIf { relative_depth: 2 }
// 0x0056 | 0e 03 00 01 |       BrTable { table: BrTable { count: 3, default: 3, targets: [0, 1, 2] } }
//        | 02 03
// 0x005c | 0f          |       Return
// 0x005d | 0b          |     End
// 0x005e | 0b          |   End
// 0x005f | 0b          | End
// 0x0060 | 0b          | End

type BrTableArgs struct {
	Labels  []LabelIdx
	Default LabelIdx
}

// ---------------- 函数调用指令
//
// call:			0x10 + func_idx
// call_indirect:	0x11 + type_idx + table_idx ;; table_idx 暂时只能是 0x00
//
// (module
//     (type $ft1 (func))
//     (type $ft2 (func))
//     (table funcref (elem $f1 $f1 $f1))
//     (func $f1
//         (call $f1)
//         (call_indirect (type $ft2) (i32.const 2))
//     )
// )
//
// 0x002a | 09          | size of function
// 0x002b | 00          | 0 local blocks
// 0x002c | 10 00       | Call { function_index: 0 }
// 0x002e | 41 02       | I32Const { value: 2 }
// 0x0030 | 11 01 00    | CallIndirect { index: 1, table_index: 0, table_byte: 0 }
// 0x0033 | 0b          | End

func (instr Instruction) GetOpname() string {
	return opnames[instr.Opcode]
}
func (instr Instruction) String() string {
	return opnames[instr.Opcode]
}
