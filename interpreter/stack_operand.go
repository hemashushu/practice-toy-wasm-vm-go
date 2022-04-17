package interpreter

import "math"

// 操作数栈（运算栈）
type operandStack struct {
	slots []uint64
}

// 部分指令是明确注明是将整数解析为有符号数再进行运算，
// 比如 lt_u, lt_s，所以需要将整数以符号数来压入和弹出的操作

// -------- 压入

func (s *operandStack) pushU64(val uint64) {
	s.slots = append(s.slots, val)
}

func (s *operandStack) pushU32(val uint32) {
	s.pushU64(uint64(val))
}

func (s *operandStack) pushS64(val int64) {
	s.pushU64(uint64(val))
}

func (s *operandStack) pushS32(val int32) {
	s.pushU32(uint32(val))
}

func (s *operandStack) pushF32(val float32) {
	s.pushU32(math.Float32bits(val))
}

func (s *operandStack) pushF64(val float64) {
	s.pushU64(math.Float64bits(val))
}

func (s *operandStack) pushBool(val bool) {
	// 使用 int32（有符号） 作为 boolean 型的数据类型
	if val {
		s.pushS32(1)
	} else {
		s.pushS32(0)
	}
}

// -------- 弹出

func (s *operandStack) popU64() uint64 {
	lastIdx := len(s.slots) - 1
	val := s.slots[lastIdx]
	s.slots = s.slots[:lastIdx]
	return val
}

func (s *operandStack) popS64() int64 {
	return int64(s.popU64())
}

func (s *operandStack) popU32() uint32 {
	return uint32(s.popU64())
}

func (s *operandStack) popS32() int32 {
	return int32(s.popU32())
}

func (s *operandStack) popF32() float32 {
	return math.Float32frombits(s.popU32())
}

func (s *operandStack) popF64() float64 {
	return math.Float64frombits(s.popU64())
}

func (s *operandStack) popBool() bool {
	// 使用 int32（有符号） 作为 boolean 型的数据类型
	return s.popS32() != 0
}

// -------- 用于实现调用栈的方法

// 栈的总大小，当前函数调用帧的实现方法是，把所有调用栈都由同一个栈来实现，所以
// sp, bp 等数值都是针对整个栈而言的
func (s *operandStack) stackSize() int {
	return len(s.slots)
}

// 按索引来获取栈的操作数
// 用于函数调用的实参以及局部变量的读写
func (s *operandStack) getOperand(idx uint32) uint64 {
	return s.slots[idx]
}

// 按索引来设置栈的操作数
// 用于函数调用的实参以及局部变量的读写
func (s *operandStack) setOperand(idx uint32, val uint64) {
	s.slots[idx] = val
}

func (s *operandStack) pushValues(vals []uint64) {
	s.slots = append(s.slots, vals...)
}

func (s *operandStack) popValues(count int) []uint64 {
	pos := len(s.slots) - count
	vals := s.slots[pos:]
	s.slots = s.slots[:pos]
	return vals
}

func (s *operandStack) peekValue() uint64 {
	lastIdx := len(s.slots) - 1
	val := s.slots[lastIdx]
	return val
}
