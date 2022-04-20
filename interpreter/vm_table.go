package interpreter

import (
	"errors"
	"wasmvm/binary"
	"wasmvm/instance"
)

// 表段和元素段目前用于列出一组函数，然后在执行 `call_indirect` 指令时，根据栈顶
// 的操作数获取该列表中的其中一个函数，从而实现 `动态` 选择被调用的函数。
// 相对应高级语言里的 `函数指针`（数据类型为 `函数` 的参数）
//
// 其中表段仅用于说明索引的大小，
// 元素段用于存储表段的初始化数据，也就是函数的索引。

// 指令 call_indirect 的操作步骤：
// 1. 从操作数栈弹出一个 uint32 数，该数是表项的索引
// 2. 从表里获取表项，从表项的值获取目标函数的索引
// 3. 通过函数索引获取目标函数
// 4. 调用目标函数
//
// |  操作数栈  |        |     表      |     | 函数列表 |
// | --------- |        | ----------- |      | ------ |
// | -- 栈顶 -- |        |  #0 func#2  |  /---> #3 sub |
// | 1:uint32 ---- pop ---> #1 func#3 ---/    | #4 mul |
// | ...       |        |  #2 func#5  |      | #5 div |
// | -- 栈底 -- |        |  ...        |      | ...    |
//
// 注
// `表` 可以被导出导入，一张表的内容，即函数引用有可能来自多个不同的模块。

type table struct {
	type_ binary.TableType // TableType 的信息包含表的类型（目前只有函数引用类型）以及限制值
	// elems []vmFunc
	elems []instance.Function
}

func NewTable(min uint32, max uint32) instance.Table {
	tableType := binary.TableType{
		ElemType: binary.FuncRef,
		Limits:   binary.Limits{Min: min, Max: max},
	}
	return newTable(tableType)
}

func newTable(tableType binary.TableType) *table {
	return &table{
		type_: tableType,
		elems: make([]instance.Function, tableType.Limits.Min),
	}
}

func (t *table) Type() binary.TableType {
	return t.type_
}

func (t *table) Size() uint32 {
	return uint32(len(t.elems))
}

func (t *table) Grow(increaseCount uint32) {
	// 检查是否超出指定的最大值
	if t.type_.Limits.Max > 0 {
		if t.Size()+increaseCount > t.type_.Limits.Max {
			panic(errors.New("out of table range"))
		}
	}

	t.elems = append(t.elems, make([]instance.Function, increaseCount)...)
}

func (t *table) GetElem(idx uint32) instance.Function {
	t.checkIdx(idx)
	elem := t.elems[idx]
	return elem
}

func (t *table) SetElem(idx uint32, elem instance.Function) {
	t.checkIdx(idx)
	t.elems[idx] = elem
}

func (t *table) checkIdx(idx uint32) {
	if idx >= uint32(len(t.elems)) {
		panic(errors.New("out of element range"))
	}
}
