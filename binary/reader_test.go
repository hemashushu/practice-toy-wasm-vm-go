package binary

import (
	"os"
	"path/filepath"
	"testing"
	"wasmvm/assert"
)

func TestReadPrimitiveData(t *testing.T) {
	reader := wasmReader{data: []byte{
		0x01,                   // byte
		0x02, 0x03, 0x04, 0x05, // 固定长度 uint32
		0x00, 0x00, 0xc0, 0x3f, // 固定长度 float32
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x3f, // 固定长度 float64
		0xE5, 0x8E, 0x26, // leb128 uint32, https://en.wikipedia.org/wiki/LEB128#Unsigned_LEB128
		0xC0, 0xBB, 0x78, // leb128 int32, https://en.wikipedia.org/wiki/LEB128#Signed_LEB128
		0xC0, 0xBB, 0x78, // leb128 int64,
		0x03, 0x01, 0x02, 0x03, // bytes/string
		0x03, 0x66, 0x6f, 0x6f, // bytes/string
	}}
	assert.AssertEqual(t, byte(0x01), reader.readByte())
	assert.AssertEqual(t, uint32(0x05040302), reader.readU32())
	assert.AssertEqual(t, float32(1.5), reader.readF32())
	assert.AssertEqual(t, 1.5, reader.readF64())
	assert.AssertEqual(t, uint32(624485), reader.readVarU32())
	assert.AssertEqual(t, int32(-123456), reader.readVarS32())
	assert.AssertEqual(t, int64(-123456), reader.readVarS64())
	assert.AssertSliceEqual(t, []byte{0x01, 0x02, 0x03}, reader.readBytes())
	assert.AssertEqual(t, "foo", reader.readName())
	assert.AssertEqual(t, 0, reader.remaining())
}

func TestReadFunctionSection(t *testing.T) {
	currentDir, err := os.Getwd() // Getwd() 返回当前 package 的目录，比如 `/path/to/project/binary`
	if err != nil {
		panic(err)
	}

	testResourcesDir := filepath.Join(currentDir, "..", "test", "resources")
	wasmFilePath := filepath.Join(testResourcesDir, "test-read-section-1.wasm")

	// "test-read-section-1.wasm" 只有一个函数和一个导出项
	//
	// (module
	// 	(export "hello" (func $hello))
	// 	(func $hello (; 0 ;) (result i32)
	// 	  (i32.const 100)
	// 	)
	// )

	m := DecodeFile(wasmFilePath)

	// 检查幻数和版本
	// 0x0000 | 00 61 73 6d | version 1 (Module)
	//        | 01 00 00 00
	assert.AssertEqual(t, MagicNumber, m.Magic)
	assert.AssertEqual(t, Version, m.Version)

	// 检查空的段
	assert.AssertEqual(t, 0, len(m.ImportSec))
	assert.AssertEqual(t, 0, len(m.TableSec))
	assert.AssertEqual(t, 0, len(m.MemSec))
	assert.AssertEqual(t, 0, len(m.GlobalSec))
	assert.AssertNil(t, m.StartSec)
	assert.AssertEqual(t, 0, len(m.ElemSec))
	assert.AssertEqual(t, 0, len(m.DataSec))

	// 检查 类型段
	// 0x0008 | 01 05       | type section
	// 0x000a | 01          | 1 count
	// 0x000b | 60 00 01 7f | [type 0] Func(FuncType { params: [], returns: [I32] })

	typeItems := m.TypeSec
	assert.AssertEqual(t, 1, len(typeItems))

	typeItem := typeItems[0]
	assert.AssertEqual(t, FtTag, typeItem.Tag)
	assert.AssertEqual(t, 0, len(typeItem.ParamTypes))
	assert.AssertSliceEqual(t, []ValType{ValTypeI32}, typeItem.ResultTypes)

	// 检查 函数（列表）段
	// 0x000f | 03 02       | func section
	// 0x0011 | 01          | 1 count
	// 0x0012 | 00          | [func 0] type 0

	funcItems := m.FuncSec
	assert.AssertEqual(t, 1, len(funcItems))
	assert.AssertSliceEqual(t, []TypeIdx{0}, funcItems)

	// 检查 导出段
	// 0x0013 | 07 09       | export section
	// 0x0015 | 01          | 1 count
	// 0x0016 | 05 68 65 6c | export Export { name: "hello", kind: Func, index: 0 }
	// 	      | 6c 6f 00 00

	exportItems := m.ExportSec
	assert.AssertEqual(t, 1, len(exportItems))

	exportItem := exportItems[0]
	assert.AssertEqual(t, "hello", exportItem.Name)
	assert.AssertEqual(t, ExportTagFunc, exportItem.Desc.Tag)
	assert.AssertEqual(t, 0, exportItem.Desc.Idx)

	// 检查 代码段
	// 0x001e | 0a 07       | code section
	// 0x0020 | 01          | 1 count
	// ============== func 0 ====================
	// 0x0021 | 05          | size of function
	// 0x0022 | 00          | 0 local blocks
	// 0x0023 | 41 e4 00    | I32Const { value: 100 }
	// 0x0026 | 0b          | End

	codeItems := m.CodeSec
	assert.AssertEqual(t, 1, len(codeItems))

	codeItem := codeItems[0]

	localsItems := codeItem.Locals
	assert.AssertEqual(t, 0, len(localsItems))

	// 代码的测试留到以后

	// 自定义段内容不检查
}

func TestReadMultipleSections(t *testing.T) {
	currentDir, err := os.Getwd() // Getwd() 返回当前 package 的目录，比如 `/path/to/project/binary`
	if err != nil {
		panic(err)
	}

	testResourcesDir := filepath.Join(currentDir, "..", "test", "resources")
	wasmFilePath := filepath.Join(testResourcesDir, "test-read-section-2.wasm")

	// "test-read-section-2.wasm" 有：
	// - 3 个类型
	// - 4 个函数
	// - 7 个导出项
	// - 3 个全局项
	// - 1 个内存块项
	//
	// (module
	// 	(type (;0;) (func (param i32 i32) (result i32)))
	// 	(type (;1;) (func (param i32) (result i32)))
	// 	(type (;2;) (func))
	// 	(func $add (type 0) (param i32 i32) (result i32)
	// 	  local.get 1
	// 	  local.get 0
	// 	  i32.add
	// 	)
	// 	(func $sub (type 0) (param i32 i32) (result i32)
	// 	  local.get 0
	// 	  local.get 1
	// 	  i32.sub
	// 	)
	// 	(func $inc (type 1) (param i32) (result i32)
	// 	  local.get 0
	// 	  i32.const 1
	// 	  i32.add
	// 	)
	// 	(func $show (type 2))
	// 	(memory (;0;) 16)
	// 	(global $__stack_pointer (mut i32) i32.const 1048576)
	// 	(global (;1;) i32 i32.const 1048576)
	// 	(global (;2;) i32 i32.const 1048576)
	// 	(export "memory" (memory 0))
	// 	(export "add" (func $add))
	// 	(export "sub" (func $sub))
	// 	(export "inc" (func $inc))
	// 	(export "show" (func $show))
	// 	(export "__data_end" (global 1))
	// 	(export "__heap_base" (global 2))
	// )

	m := DecodeFile(wasmFilePath)

	// 检查幻数和版本
	assert.AssertEqual(t, MagicNumber, m.Magic)
	assert.AssertEqual(t, Version, m.Version)

	// 检查空的段
	assert.AssertEqual(t, 0, len(m.ImportSec))
	assert.AssertEqual(t, 0, len(m.TableSec))
	assert.AssertNil(t, m.StartSec)
	assert.AssertEqual(t, 0, len(m.ElemSec))
	assert.AssertEqual(t, 0, len(m.DataSec))

	// 检查 类型段
	// [type 0] Func(FuncType { params: [I32, I32], returns: [I32] })
	// [type 1] Func(FuncType { params: [I32], returns: [I32] })
	// [type 2] Func(FuncType { params: [], returns: [] })

	typeItems := m.TypeSec
	assert.AssertEqual(t, 3, len(typeItems))

	typeItem0 := typeItems[0]
	assert.AssertEqual(t, FtTag, typeItem0.Tag)
	assert.AssertSliceEqual(t, []ValType{ValTypeI32, ValTypeI32}, typeItem0.ParamTypes)
	assert.AssertSliceEqual(t, []ValType{ValTypeI32}, typeItem0.ResultTypes)

	typeItem1 := typeItems[1]
	assert.AssertEqual(t, FtTag, typeItem1.Tag)
	assert.AssertSliceEqual(t, []ValType{ValTypeI32}, typeItem1.ParamTypes)
	assert.AssertSliceEqual(t, []ValType{ValTypeI32}, typeItem1.ResultTypes)

	typeItem2 := typeItems[2]
	assert.AssertEqual(t, FtTag, typeItem2.Tag)
	assert.AssertSliceEqual(t, []ValType{}, typeItem2.ParamTypes)
	assert.AssertSliceEqual(t, []ValType{}, typeItem2.ResultTypes)

	// 检查 函数（列表）段
	// [func 0] type 0
	// [func 1] type 0
	// [func 2] type 1
	// [func 3] type 2

	funcItems := m.FuncSec
	assert.AssertEqual(t, 4, len(funcItems))
	assert.AssertSliceEqual(t, []TypeIdx{0, 0, 1, 2}, funcItems)

	// 检查 内存段
	// [memory 0] MemoryType { memory64: false, shared: false, initial: 16, maximum: None }
	memBlockItems := m.MemSec
	assert.AssertEqual(t, 1, len(memBlockItems))

	memBlockItem := memBlockItems[0]
	assert.AssertEqual(t, 0, memBlockItem.Tag)
	assert.AssertEqual(t, 16, memBlockItem.Min)
	assert.AssertEqual(t, 0, memBlockItem.Max)

	// 检查 全局段
	// [global 0] GlobalType { content_type: I32, mutable: true }
	// [global 1] GlobalType { content_type: I32, mutable: false }
	// [global 2] GlobalType { content_type: I32, mutable: false }
	globalItems := m.GlobalSec
	assert.AssertEqual(t, 3, len(globalItems))

	globalItem0 := globalItems[0]
	assert.AssertEqual(t, ValTypeI32, globalItem0.Type.ValType)
	assert.AssertEqual(t, MutVar, globalItem0.Type.Mut)
	// 代码的测试留到以后

	globalItem1 := globalItems[1]
	assert.AssertEqual(t, ValTypeI32, globalItem1.Type.ValType)
	assert.AssertEqual(t, MutConst, globalItem1.Type.Mut)
	// 代码的测试留到以后

	globalItem2 := globalItems[2]
	assert.AssertEqual(t, ValTypeI32, globalItem2.Type.ValType)
	assert.AssertEqual(t, MutConst, globalItem2.Type.Mut)
	// 代码的测试留到以后

	// 检查 导出段
	// export Export { name: "memory", kind: Memory, index: 0 }
	// export Export { name: "add", kind: Func, index: 0 }
	// export Export { name: "sub", kind: Func, index: 1 }
	// export Export { name: "inc", kind: Func, index: 2 }
	// export Export { name: "show", kind: Func, index: 3 }
	// export Export { name: "__data_end", kind: Global, index: 1 }
	// export Export { name: "__heap_base", kind: Global, index: 2 }

	exportItems := m.ExportSec
	assert.AssertEqual(t, 7, len(exportItems))

	exportItem0 := exportItems[0]
	assert.AssertEqual(t, "memory", exportItem0.Name)
	assert.AssertEqual(t, ExportTagMem, exportItem0.Desc.Tag)
	assert.AssertEqual(t, 0, exportItem0.Desc.Idx)

	exportItem1 := exportItems[1]
	assert.AssertEqual(t, "add", exportItem1.Name)
	assert.AssertEqual(t, ExportTagFunc, exportItem1.Desc.Tag)
	assert.AssertEqual(t, 0, exportItem1.Desc.Idx)

	exportItem2 := exportItems[2]
	assert.AssertEqual(t, "sub", exportItem2.Name)
	assert.AssertEqual(t, ExportTagFunc, exportItem2.Desc.Tag)
	assert.AssertEqual(t, 1, exportItem2.Desc.Idx)

	exportItem3 := exportItems[3]
	assert.AssertEqual(t, "inc", exportItem3.Name)
	assert.AssertEqual(t, ExportTagFunc, exportItem3.Desc.Tag)
	assert.AssertEqual(t, 2, exportItem3.Desc.Idx)

	exportItem4 := exportItems[4]
	assert.AssertEqual(t, "show", exportItem4.Name)
	assert.AssertEqual(t, ExportTagFunc, exportItem4.Desc.Tag)
	assert.AssertEqual(t, 3, exportItem4.Desc.Idx)

	exportItem5 := exportItems[5]
	assert.AssertEqual(t, "__data_end", exportItem5.Name)
	assert.AssertEqual(t, ExportTagGlobal, exportItem5.Desc.Tag)
	assert.AssertEqual(t, 1, exportItem5.Desc.Idx)

	exportItem6 := exportItems[6]
	assert.AssertEqual(t, "__heap_base", exportItem6.Name)
	assert.AssertEqual(t, ExportTagGlobal, exportItem6.Desc.Tag)
	assert.AssertEqual(t, 2, exportItem6.Desc.Idx)

	// 检查 代码段
	codeItems := m.CodeSec
	assert.AssertEqual(t, 4, len(codeItems))
	// 代码的测试留到以后

	// 自定义段内容不检查
}
