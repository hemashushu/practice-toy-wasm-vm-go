package binary

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
)

type wasmReader struct {
	data []byte // 待读取的二进制数据
}

func DecodeFile(filename string) Module {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return Decode(data)
}

func Decode(data []byte) Module {
	module := Module{}
	reader := &wasmReader{data: data}
	reader.readModule(&module)
	return module
}

// -------- 辅助读取函数

// 读取一个字节
func (r *wasmReader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

// 读取固定长度的 uint32
func (r *wasmReader) readU32() uint32 {
	n := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:] // 消耗掉 4 bytes
	return n
}

// 读取固定长度的 float32
func (r *wasmReader) readF32() float32 {
	n := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:] // 消耗掉 4 bytes
	return math.Float32frombits(n)
}

// 读取固定长度的 float64
func (r *wasmReader) readF64() float64 {
	n := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:] // 消耗掉 8 bytes
	return math.Float64frombits(n)
}

// 读取变长（leb128）uint32
func (r *wasmReader) readVarU32() uint32 {
	value, bytes := decodeVarUint(r.data, 32)
	r.data = r.data[bytes:]
	return uint32(value)
}

// 读取变长（leb128）signed int32
func (r *wasmReader) readVarS32() int32 {
	value, bytes := decodeVarInt(r.data, 32)
	r.data = r.data[bytes:]
	return int32(value)
}

// 读取变长（leb128）signed int64
func (r *wasmReader) readVarS64() int64 {
	value, bytes := decodeVarInt(r.data, 64)
	r.data = r.data[bytes:]
	return value
}

// 读取字节数组
// 项目内容的长度位于开头的一个 uint32 数字
func (r *wasmReader) readBytes() []byte {
	length := r.readVarU32()
	bytes := r.data[:length]
	r.data = r.data[length:]
	return bytes
}

// 读取字符串
// 字符串的长度位于开头的一个 uint32 数字
func (r *wasmReader) readName() string {
	data := r.readBytes()
	return string(data)
}

// 获取剩余的数据的长度（字节数）
func (r *wasmReader) remaining() int {
	return len(r.data)
}

// -------- 解码模块

func (r *wasmReader) readModule(m *Module) {
	m.Magic = r.readU32()
	m.Version = r.readU32()

	r.readSections(m)
}

func (r *wasmReader) readSections(m *Module) {
	// 记录上一次/最后一次解析段 id，用于确保段 id
	// 是按照正确顺序出现
	lastSectionId := byte(0)

	for r.remaining() > 0 {
		sectionId := r.readByte()
		if sectionId == SecCustomID {
			// 自定义段的出现顺序不固定，而且可能出现多次
			m.CustomSecs = append(m.CustomSecs,
				r.readCustomSec())
		} else {
			// 除了自定义段，其他段的 id 出现顺序是按照从小到大的顺序出现
			if sectionId > SecDataID || sectionId <= lastSectionId {
				panic(errors.New("invalid section id"))
			}
			lastSectionId = sectionId

			// 当前段的长度，注意这里已经消耗了当前段的长度数据
			length := r.readVarU32()
			// 当前剩余的长度
			lastRemainBytes := r.remaining()

			r.readNonCustomSec(sectionId, m)

			// 检查段解析过程是否正确地解析完当前段的所有数据
			remainBytes := r.remaining()
			if remainBytes != lastRemainBytes-int(length) {
				panic(errors.New("section parser consumed unexpected length of data"))
			}
		}
	}
}

// -------- 解码自定义段

func (r *wasmReader) readCustomSec() CustomSec {
	// 自定义段的数据长度由段 id 后的第一个 uint32 数字指出

	// 消耗整个 custom section 的有效数据，然后构造新的解析器
	sectionReader := wasmReader{data: r.readBytes()}
	return CustomSec{
		Name:  sectionReader.readName(),
		Bytes: sectionReader.data,
	}
}

// -------- 解码非自定义段（的入口）

func (r *wasmReader) readNonCustomSec(sectionId byte, m *Module) {
	switch sectionId {
	case SecTypeID:
		m.TypeSec = r.readTypeSec()
	case SecImportID:
		m.ImportSec = r.readImportSec()
	case SecFuncID:
		m.FuncSec = r.readFuncSec()
	case SecTableID:
		m.TableSec = r.readTableSec()
	case SecMemID:
		m.MemSec = r.readMemSec()
	case SecGlobalID:
		m.GlobalSec = r.readGlobalSec()
	case SecExportID:
		m.ExportSec = r.readExportSec()
	case SecStartID:
		m.StartSec = r.readStartSec()
	case SecElemID:
		m.ElemSec = r.readElemSec()
	case SecCodeID:
		m.CodeSec = r.readCodeSec()
	case SecDataID:
		m.DataSec = r.readDataSec()
	}
}

// -------- 解码（函数）类型段

// 段的内容长度位于第一个字节之后的一个 uint32 数字
// 第一个字节是段的类型 id

func (r *wasmReader) readTypeSec() []FuncType {
	vec := make([]FuncType, r.readVarU32())
	for i := range vec {
		vec[i] = r.readFuncType()
	}
	return vec
}

func (r *wasmReader) readFuncType() FuncType {
	ft := FuncType{
		Tag:         r.readByte(),
		ParamTypes:  r.readValTypes(),
		ResultTypes: r.readValTypes(),
	}

	if ft.Tag != FtTag {
		panic(fmt.Errorf("invalid function type tag: %d", ft.Tag))
	}

	return ft
}

func (r *wasmReader) readValTypes() []ValType {
	vec := make([]ValType, r.readVarU32())
	for i := range vec {
		vec[i] = r.readValType()
	}
	return vec
}

// 因为数据类型只有 4 种，所以 data_type 的类型是 byte
func (reader *wasmReader) readValType() ValType {
	b := reader.readByte()

	if b != ValTypeI32 &&
		b != ValTypeI64 &&
		b != ValTypeF32 &&
		b != ValTypeF64 {
		panic(fmt.Errorf("invalid data type: %d", b))
	}

	return b
}

// -------- 解码导入段

func (r *wasmReader) readImportSec() []Import {
	vec := make([]Import, r.readVarU32())
	for i := range vec {
		vec[i] = r.readImport()
	}
	return vec
}

func (r *wasmReader) readImport() Import {
	return Import{
		Module: r.readName(),
		Name:   r.readName(),
		Desc:   r.readImportDesc(),
	}
}

func (r *wasmReader) readImportDesc() ImportDesc {
	desc := ImportDesc{Tag: r.readByte()}
	switch desc.Tag {
	case ImportTagFunc:
		// 函数索引
		desc.FuncType = r.readVarU32()
	case ImportTagTable:
		// 表项目
		desc.Table = r.readTableType()
	case ImportTagMem:
		// 内存块项目
		desc.Mem = r.readLimits()
	case ImportTagGlobal:
		// 全局变量项目
		desc.Global = r.readGlobalType()
	default:
		panic(fmt.Errorf("invalid import desc tag: %d", desc.Tag))
	}
	return desc
}

// -------- 解码函数（列表）段

func (r *wasmReader) readFuncSec() []TypeIdx {
	vec := make([]TypeIdx, r.readVarU32())
	for i := range vec {
		vec[i] = r.readVarU32()
	}
	return vec
}

// -------- 解码表（列表）段

func (r *wasmReader) readTableSec() []TableType {
	vec := make([]TableType, r.readVarU32())
	for i := range vec {
		vec[i] = r.readTableType()
	}
	return vec
}

func (r *wasmReader) readTableType() TableType {
	tt := TableType{
		ElemType: r.readByte(),
		Limits:   r.readLimits(),
	}
	if tt.ElemType != FuncRef {
		panic(fmt.Errorf("invalid elemtype: %d", tt.ElemType))
	}
	return tt
}

func (r *wasmReader) readLimits() Limits {
	limits := Limits{
		Tag: r.readByte(),
		Min: r.readVarU32(),
	}

	// 仅当 tag == 1 时，才有 max 数据
	if limits.Tag == 1 {
		limits.Max = r.readVarU32()
	}
	return limits
}

// -------- 解码内存块（列表）段

func (r *wasmReader) readMemSec() []MemType {
	vec := make([]MemType, r.readVarU32())
	for i := range vec {
		vec[i] = r.readLimits()
	}
	return vec
}

// -------- 解码全局变量段

func (r *wasmReader) readGlobalSec() []Global {
	vec := make([]Global, r.readVarU32())
	for i := range vec {
		vec[i] = Global{
			Type: r.readGlobalType(),
			Init: r.readExpr(), // 初始值表达式
		}
	}
	return vec
}

func (r *wasmReader) readGlobalType() GlobalType {
	gt := GlobalType{
		ValType: r.readValType(),
		Mut:     r.readByte(),
	}
	if gt.Mut != MutConst &&
		gt.Mut != MutVar {
		panic(errors.New("invalid global mutability type"))
	}
	return gt
}

func (r *wasmReader) readExpr() Expr {
	for r.readByte() != 0x0B {
		// todo::
		// read until reach 0x0B
	}

	return nil
}

// -------- 解码导出段

func (r *wasmReader) readExportSec() []Export {
	vec := make([]Export, r.readVarU32())
	for i := range vec {
		vec[i] = r.readExport()
	}
	return vec
}

func (reader *wasmReader) readExport() Export {
	return Export{
		Name: reader.readName(),
		Desc: reader.readExportDesc(),
	}
}

func (reader *wasmReader) readExportDesc() ExportDesc {
	desc := ExportDesc{
		Tag: reader.readByte(),
		Idx: reader.readVarU32(),
	}

	if desc.Tag != ExportTagFunc && // func_idx
		desc.Tag != ExportTagTable && // table_idx
		desc.Tag != ExportTagMem && // mem_idx
		desc.Tag != ExportTagGlobal { // global_idx
		panic(errors.New("invalid export desc tag"))
	}
	return desc
}

// -------- 解码起始段

func (r *wasmReader) readStartSec() *uint32 {
	idx := r.readVarU32()
	return &idx
}

// -------- 解码元素段（表的初始内容）

func (r *wasmReader) readElemSec() []Elem {
	vec := make([]Elem, r.readVarU32())
	for i := range vec {
		vec[i] = r.readElem()
	}
	return vec
}

func (r *wasmReader) readElem() Elem {
	return Elem{
		Table:  r.readVarU32(),      // 表索引
		Offset: r.readExpr(),        // 偏移值的表达式
		Init:   r.readFuncIndices(), // 函数索引
	}
}

func (r *wasmReader) readFuncIndices() []FuncIdx {
	vec := make([]FuncIdx, r.readVarU32())
	for i := range vec {
		vec[i] = r.readVarU32()
	}
	return vec
}

// -------- 解码（函数）代码段

func (r *wasmReader) readCodeSec() []Code {
	vec := make([]Code, r.readVarU32())
	for i := range vec {
		vec[i] = r.readCode()
	}
	return vec
}

func (r *wasmReader) readCode() Code {
	codeReader := wasmReader{data: r.readBytes()}
	code := Code{
		Locals: codeReader.readLocalsVec(),
		// todo::
		//Expr:   reader.readExpr(),
	}

	// 检查局部变量的数量是否溢出
	total := uint64(0)
	for _, locals := range code.Locals {
		total += uint64(locals.N)
	}

	if total >= math.MaxUint32 {
		panic(fmt.Errorf("too many locals: %d", total))
	}

	return code
}

func (r *wasmReader) readLocalsVec() []Locals {
	vec := make([]Locals, r.readVarU32())
	for i := range vec {
		vec[i] = r.readLocals()
	}
	return vec
}

func (r *wasmReader) readLocals() Locals {
	return Locals{
		N:    r.readVarU32(),
		Type: r.readValType(),
	}
}

// -------- 解码数据段（内存块初始内容）

func (r *wasmReader) readDataSec() []Data {
	vec := make([]Data, r.readVarU32())
	for i := range vec {
		vec[i] = r.readData()
	}
	return vec
}

func (r *wasmReader) readData() Data {
	return Data{
		Mem:    r.readVarU32(),
		Offset: r.readExpr(),
		Init:   r.readBytes(),
	}
}
