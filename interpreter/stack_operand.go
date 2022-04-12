package interpreter

import "math"

// 操作数栈（运算栈）
type operandStack struct {
	slots []uint64
}

// -------- 压入

func (s *operandStack) pushU64(val uint64) {
	s.slots = append(s.slots, val)
}

func (s *operandStack) pushS64(val int64) {
	s.pushU64(uint64(val))
}

func (s *operandStack) pushU32(val uint32) {
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
