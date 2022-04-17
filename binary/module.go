package binary

// module:
//   magic:uint32 + version:uint32 +
//   type_sec? +
//   import_sec? +
//   func_sec? +
//   table_sec? +
//   mem_sec? +
//   global_sec? +
//   export_sec? +
//   start_sec? +
//   elem_sec? +
//   code_sec? +
//   data_sec?

type Module struct {
	Magic   uint32 // 幻数
	Version uint32 // 版本号

	//-----------------------     V-- 下面这些数字是段 id
	CustomSecs []CustomSec // 编号 0: 自定义段
	TypeSec    []FuncType  // 编号 1: 类型段（即函数签名列表，不同的函数可能有相同的签名，所以这里的列表是排除重复了之后的）
	ImportSec  []Import    // 编号 2: 导入（函数信息）段
	FuncSec    []TypeIdx   // 编号 3: 函数列表段，列出所有的函数（注意不同的函数可能有相同的签名）的列表，跟代码段合在一起形成完整的函数
	TableSec   []TableType // 编号 4: 表格段，表格段跟元素段合在一起实现函数间接调用，!! 目前只支持 1 项。
	MemSec     []MemType   // 编号 5: 内存（描述）段，!! 目前只支持 1 项。
	GlobalSec  []Global    // 编号 6: 全局变量信息段
	ExportSec  []Export    // 编号 7: 导出（函数信息）段
	StartSec   *FuncIdx    // 编号 8: 程序入口函数（即主函数、main 函数）
	ElemSec    []Elem      // 编号 9: 元素段，跟表格段合在一起实现函数间接调用
	CodeSec    []Code      // 编号 10: 函数主体段，跟函数列表段合在一起实现完整的函数
	DataSec    []Data      // 编号 11: （内存初始）数据段，跟内存描述段合在一起形成完整的初始数据
}

const (
	SecCustomID = iota // 0
	SecTypeID          // 1
	SecImportID        // 2
	SecFuncID          // 3
	SecTableID         // 4
	SecMemID           // 5
	SecGlobalID        // 6
	SecExportID        // 7
	SecStartID         // 8
	SecElemID          // 9
	SecCodeID          // 10
	SecDataID          // 11
)

// 类型别名（仅为了提高代码可读性）

type (
	TypeIdx   = uint32 // 函数类型索引（内部、导入函数共用）
	FuncIdx   = uint32 // 函数索引（内部、导入函数共用）
	GlobalIdx = uint32 // 全局变量索引
	TableIdx  = uint32 // 表索引，目前的值只能是 0
	MemIdx    = uint32 // 内存索引，目前的值只能是 0
	LocalIdx  = uint32 // （每个函数的）局部变量索引
	LabelIdx  = uint32 // （每个函数内部）跳转标签的索引
)

// 二进制文件头

const (
	MagicNumber = 0x6d736100 // "0asm", 占用了 4 个字节，在内存中应该是  "msa0"
	Version     = 0x00000001 // 1， 占用了 4 个字节
)

// ================ 段的定义

// ---------------- 段的前缀

// 每一个段的开头
// id:byte + byte_count:uint32 + ...
//
// 其中 byte_count 是指该段的正文内容部分的长度（即不包括头部自己）

// 段头
type SecHeader struct {
	Id    byte   // 段的 id
	Bytes uint32 // 段内容的长度
}

// ---------------- （函数）类型段（函数类型即函数签名）

// type_sec: 0x01 + byte_count:byte + func_type_items_count:uint32 + func_type{1,}
//
// 以下使用 "<...>" 代表 "some_type_items_count:uint32 + some_type{1,}" 这种结构，
// 比如
// "<func_type>" == "func_type_items_count:uint32 + func_type{1,}"
//
// type_sec: 0x01 + byte_count + type_count + <func_type>
// func_type: 0x60 + <val_type> + <val_type>
//                   ^                 ^
//                   |--- 参数类型列表   |--- 返回值类型列表
//
// <val_type>: count:uint32 + data_type:byte{0,}
//
// 因为数据类型只有 4 种，所以 data_type 的数据类型是 byte
//
// 文本格式
//
// (type (func (param i32) (param i32) (result i32)))
// ;; 可以添加自动索引 $ft1（诸如 `$xx` 在 wat 里也叫标识符）
// (type $ft1 (func (param i32) (result f64)))
// ;; 多个参数可以写在同一个 param 列表里，多个返回值也可以写在同一个 result 列表里
// (type $ft1 (func (param i32 i32) (result f64 f64)))

// （函数）类型项目
type FuncType struct {
	Tag         byte      // 只能是 0x60
	ParamTypes  []ValType // 参数类型列表
	ResultTypes []ValType // 返回值类型列表
}

const FtTag = 0x60

// ---------------- 导入段

// import_sec: 0x02 + byte_count:uint32 + <import>
// import: module_name:string + member_name:string + import_desc
//
// module_name 和 member_name 是字符串以 utf-8 编码后的字节数组，在字节数组
// 之前有一个 uint32 描述字节数组的长度（这个长度是只字符串正文内容的长度，所以
// 当然不包括这个描述数字本身占用的空间）
//
// import_desc: tag:byte + (func_type_idx | table_type | mem_type | global_type)
//
// 文本格式：
//
// ;; "env" 和 "f1" 分别是导入项的模块名和条目名
// (type $ft1 (func (param i32 i32) (result i32)))
// (import "env" "f1" (func $f1 (type $ft1)))
//
// ;; $ft1 和 “f1" 可以内联
// (import "env" "f1" (func $f1 (param i32 i32) (result i32)))
//
// (import "env" "t1" (table $t 1 8 funcref))
// (import "env" "m1" (memory $m 4, 16))
// (import "env" "g1" (global $g1 i32))			;; 全局常量
// (import "env" "g2" (global $g1 (mut i32)))	;; 全局变量
//
// 导入项可以内联到函数、表、内存和全局中
//
// (func $f1 (import "env" "f1") (type $ft1))
// (table $t1 (import "env" "t1") 1 8 funcref)
// (memory $m1 (import "env" "m1") 4 16)
// (global $g1 (import "env" "g1") i32)
// (global $g2 (import "env" "g2") (mut i32))

// 导入项目
type Import struct {
	Module string     // 模块名称
	Name   string     // 项目名称（比如函数名）
	Desc   ImportDesc // 导入项的描述
}

// 导入项类型
const (
	ImportTagFunc   = 0
	ImportTagTable  = 1
	ImportTagMem    = 2
	ImportTagGlobal = 3
)

// 导入项描述
type ImportDesc struct {
	Tag      byte       // 导入项类型
	FuncType TypeIdx    // 仅当 tag == 0 时有效
	Table    TableType  // 仅当 tag == 1 时有效
	Mem      MemType    // 仅当 tag == 2 时有效
	Global   GlobalType // 仅当 tag == 3 时有效
}

// 函数（列表）段

// func_sec: 0x03 + byte_count:uint32 + <type_idx>
//
// 函数列表仅列出该函数的类型数字，比如 (func_)type_sec 里有 2 条记录：
// type0: (u32, u32) u32
// type1: (f32) u32
//
// 则函数列表 `00 01 01 00`
// 表示一共有 4 个函数，
//
// func0 (u32, u32) u32
// func1 (f32) u32
// func2 (f32) u32
// func3 (u32, u32) u32
//
// 注意
// 函数的索引有可能不是从 0 开始，比如导入了 3 个函数，则这个列表的第一个函数
// 的索引应该是 3。
//
// 文本格式
//
// (type $ft1 (func (param i32 i32) (result i32)))
// (func $add (type $ft1)	;; $ft1 是类型索引，$add 是函数的（自动）索引
//    (local i64 i64)		;; 声明两个局部变量
//    (i64.add (local.get 2) (local.get 3))	;; 访问上面两个局部变量，local.get 指令使用了内联方式书写
//    (drop)
//    (i32.add (local.get 0) (local.get 1))	;; 访问函数的两个参数，函数参数也是局部变量
// )
//
// 如果不使用内联方式书写局部变量，则：
//
// (func $add (param $a i32) (param $b i32) (result i32)
// 	(local $x i64)
// 	(local $y i64)
// 	(i64.add (local.get $x) (local.get $y))	;; 索引数字换成了自动索引（的名称）
// 	(drop)
// 	(i32.add (local.get $a) (local.get $b))
// )

// ---------------- 表段

// 表段和元素段目前可用于列出 “指针化” 的函数，实现诸如 “高阶函数” 和跨模块调用等功能。
// 其中表段仅用于说明索引的大小，真正的函数索引列表在元素段里，
// 也就是说元素段存储的是表的（初始化）数据
//
// table_sec: 0x04 + byte_count:uint32 + <table_type> // 目前仅支持一个 table_type
// table_type: 0x70 + limits
//
// (func $f1)
// (func $f2)
// (table 1 10 funcref)						;; 表的类型暂时只能是 "funcref"
// (elem (offset (i32.const 1)) $f1 $f2)	;; 元素项的偏移值需要使用（const）表达式
//
// 元素项也可以内联到表段里：
//
// (table funcref				;; 自动决定了表的 limit, min = 2, max = 2
// 	(elem $f1 $f2)				;; 自动决定了偏移值为 0
// )

const FuncRef = 0x70

// 表项目
type TableType struct {
	ElemType byte   // 表项目的类型 目前只支持 0x70
	Limits   Limits // 限制值
}

// 限制类型（用于描述元素数量/内存页数的上下限）

// limits: tag:byte + min + max（可选）
//
// min 是下限值，max 是上限值
// 当 tag == 0 时，表示省略了上限，只有 min 值
// 当 tag == 1 时，表示上下限都指出
//
// 示例：
// 00 01      ; 下限值为 1，省略了上限（所以上限的字节也不会有）
// 01 01 02   ; 下限值为 1，上限值为 2

// 限制值
type Limits struct {
	Tag byte   // 限制值的类型，0 表示只有 min 值，1 表示有 min 和 max 值
	Min uint32 // min，即下限
	Max uint32 // max，即上限，是可选的，省略上限时，该位置对应的字节也不会有
}

// ---------------- 内存段

// mem_sec: 0x05 + byte_count:uint32 + <mem_type> // 目前仅支持一个 mem_type
// mem_type: limits
//
// 文本格式
//
// (memory 1 16)						;; 指定 limit 值，即 min 和 max
// (data (offset (i32.const 10)) "foo")	;; 数据偏移量需要使用（const）表达式
// (data (offset (i32.const 20)) "bar")
//
// 将数据内联到内存段里：
// (memory 			;; 自动 limit, min = 1, max = 1
// 	(data "foo")	;; 自动偏移值 0
// 	(data "bar")
// )
//
// 初始数据使用字符串的形式指定，内容可以是
// - 单一字符："abc文字"（字符将会以 utf-8 形式编码）
// - 十六进制 byte: "\de\ad\be\ef\00"
// - Unicode code point: "\u{1234}\u{5678}"

type MemType = Limits

const (
	PageSize     = 65536 // 每页内存 64KB
	MaxPageCount = 65536
)

// ---------------- 全局段

// 全局段列出模块所有全局变量/常量
// 变量需要指出是否可变，以及初始值表达式
//
// gloabl_sec: 0x06 + byte_count:uint32 + <global>
// gloabl: gloabl_type + init_expr
// gloabl_type: val_type:byte + mut
// mut: 0:byte（不可变） | 1:byte（可变）
// init_expr: byte{0,} + 0x0B // 初始值表达式以 0x0B 结尾
//
// 全局项示例：
// - 7f             ; 当前全局变量的数据类型是 i32
// - 01   			; 可变
// - 41 80 80 c0 00	; i32.const 0x100000
// - 0b				; 初始值表达式结束
//
// 文本格式
// (gloabl $g1 (mut i32) (i32.const 10))	;; 全局变量
// (gloabl $g1 i32 (i32.const 20))			;; 全局常量
// (func
// 	(global.get $g1)
// 	(global.get $g2)
// )

// 全局项目
type Global struct {
	Type GlobalType // 全局项目属性
	Init Expr       // 初始值表达式（指令/字节码）
}

// 全局项目的属性
type GlobalType struct {
	ValType ValType // 数据类型
	Mut     byte    // 是否可变
}

// 全局变量类型
const (
	MutConst byte = 0 // 常量（这里名称前缀 `Mut` 起得不太好，应该叫做 GlobalConst/GlobalVar 更合适）
	MutVar   byte = 1 // 变量
)

// ---------------- 导出段

// 可以导出：函数、表、内存、全局变量
// export_sec: 0x07 + byte_count:uint32 + <export>
// export: name:string + export_desc
// export_desc: tag:byte + (func_idx | table_idx | mem_idx | global_idx)
//
// 文本格式：
//
// (export "f1" (func $f1))		;; 使用自动索引 #f1
// (export "f2" (func $f2))
// (export "t1" (table $t))
// (export "m1" (memory $m))
// (export "g1" (global $g1))
// (export "g2" (global $g2))
//
// 导出项内联到函数、表、内存、全局
//
//	(func $f (export "f1") ...)
//	(func $t (export "t") ...)
//	(func $m (export "m") ...)
//	(func $g (export "g1") ...)

// 导出项
type Export struct {
	Name string     // 导出项的名称（导出项不需要指定当前模块的名称，导入时则需同时指出导入模块和导入项的名称）
	Desc ExportDesc // 导出项描述
}

// 导出项类型
const (
	ExportTagFunc   = 0
	ExportTagTable  = 1
	ExportTagMem    = 2
	ExportTagGlobal = 3
)

// 导出项描述
type ExportDesc struct {
	Tag byte   // 导出项类型
	Idx uint32 // 函数、表、内存块、全局项的索引
}

// ---------------- 起始段

// 指定 wasm 加载后自动开始执行的函数（比如 main 函数）
// start_sec: 0x08 + byte_count:uint32 + func_idx
//
// 文本格式
//
// (module
// 	(func $main ...)
// 	(start $main)	;; start 指令后面跟着起始函数的索引值
// )

// ---------------- 元素段

// 元素段用于存储存放表段的初始化数据，跟表段共同完成诸如 “高阶函数” 的功能，
// 目前元素段的内容是函数的索引。
//
// 元素段里的每个项目的内容由 3 部分组成：
// 1. 表的索引（目前只能是 0）
// 2. 表内偏移量（是一个表达式）
// 3. 函数索引列表
//
// elem_sec: 0x09 + byte_count:uint32 + <elem>
// elem: table_idx + offset_expr + <func_idx>

// 元素项目
type Elem struct {
	Table  TableIdx  // 表索引
	Offset Expr      // 偏移值表达式（指令/字节码）
	Init   []FuncIdx // 函数索引
}

// ---------------- 代码段

// code_sec: 0x0a + byte_count:uint32 + <code>
// code: byte_count:uint32 + <locals> + expr
// locals: local_count + val_type
//
// code 项开头的 byte_count 表示该项目的内容大小，包括表达式结尾的 0x0B
//
// 示例：
//
// (func (param $a i32) (param $b i32)
// 	(local $la i32)
// 	(local $lb i32)
// 	(local i64 i64)
// 	(global.get $g1)
// 	(global.set $g2)
// 	(local.get $a)
// 	(local.set $b)
// )
//
// - 0e          | size of function
// - 02          | 2 local blocks
// - 02 7f       | 2 locals of type I32
// - 02 7e       | 2 locals of type I64
// - 23 00       | GlobalGet { global_index: 0 }
// - 24 01       | GlobalSet { global_index: 1 }
// - 20 00       | LocalGet { local_index: 0 }
// - 21 01       | LocalSet { local_index: 1 }
// - 0b          | End

// 代码项目（一个函数一个代码项目）
type Code struct {
	Locals []Locals // 局部变量组列表，连续多个相同类型的局部变量被分为一组
	Expr   Expr     // 指令/字节码
}

// 局部变量组
type Locals struct {
	N    uint32  // 数量
	Type ValType // 类型
}

// 辅助函数
// 用于将 block 的返回类型跟函数的签名（function type）统一起来，
// 因为 block type 都是负数，所以不会跟 function type 的索引冲突
func (module Module) GetBlockType(bt BlockType) FuncType { // name: convertBlockTypeIntoFunctionType
	switch bt {
	case BlockTypeI32: // -1
		return FuncType{ResultTypes: []ValType{ValTypeI32}}
	case BlockTypeI64: // -2
		return FuncType{ResultTypes: []ValType{ValTypeI64}}
	case BlockTypeF32: // -3
		return FuncType{ResultTypes: []ValType{ValTypeF32}}
	case BlockTypeF64: // -4
		return FuncType{ResultTypes: []ValType{ValTypeF64}}
	case BlockTypeEmpty: // -64
		return FuncType{}
	default:
		return module.TypeSec[bt]
	}
}

func (code Code) GetLocalCount() uint64 {
	n := uint64(0)
	for _, locals := range code.Locals {
		n += uint64(locals.N)
	}
	return n
}

// ---------------- 数据段

// 数据段跟元素段类似，存储内存的初始化数据
// 数据段每一项由 3 部分组成：
// 1. 内存块索引（目前只能是 0）
// 2. 内存偏移值（是一个表达式）
// 3. 初始数据
//
// 数据项示例：
// - 00				; 内存块索引
// - 41 80 80 c0 00	; i32.const(0x41) 0x100000
// - 0b				; 偏移值表达式结束标记
// - 0e				; 内容长度 14 字节（0x0e）
// - 48 65 6c 6c 6f ; "Hello"
// - 2c 20 			; ", "
// - 57 6f 72 6c 64 ; "World"
// - 21 0A			; "!\n"
//
// data_sec: 0x0b + byte_count:uint32 + <data>
// data: mem_idx + offset_expr + byte{0,}

// 数据项目
type Data struct {
	Mem    MemIdx // 内存块索引
	Offset Expr   // 偏移值表达式（指令/字节码）
	Init   []byte // 内容
}

// ---------------- 自定义段
//
// 自定义段可以出现多次，出现的位置也不限。
// 一般用于存放函数的名称、参数和变量的名称等信息，不参与运算
//
// custom_sec: 0x00 + byte_count:uint32 + name:string + byte{0,}

type CustomSec struct {
	Name  string
	Bytes []byte
}
