package interpreter

import "wasmvm/binary"

// 当前的 vm 实现共用调用帧以及流程控制的块帧，所以
// 称为 `控制帧`（`controlFrame`）。
//
// 当前的 vm 实现不为每个调用帧的创建新的操作数栈，而是
// 将所有操作数栈都共享同一个操作数栈，然后使用控制帧记录
// 当前帧的开始位置（也就是第 0 个实参的位置），这样可以避免复制实参。
// 但注意的是，这种实现方法仍然无法避免返回值的复制，因为有时当前函数可能因为逻辑错误
// 导致当前操作数栈内有除了返回值之外的操作数残留。当然如果没有残留，是可以不用复制返回值。

// 逻辑栈帧示意图
//
//
//                         当前栈帧 -- | ------- 栈顶 -------- |
//                                    | 运算槽位               |
//                                    | -------------------- |
//                                    | 局部变量占用的槽位      |
//        | ------- 栈顶 -------- |    | -------------------- | <--
//        | 传给下一个栈帧的实参    | -> | 来自上一个栈帧的实参     |   |
//        |                      | -> |                      |   |-- 重叠区域
//        | .................... | -> | ------- 栈底 -------- | <--
// 上     | 运算槽位               |
// 一个    | -------------------- |
// 栈帧 -- | 局部变量占用的槽位      |
//        | -------------------- |
//        | 来自上一个栈帧的实参    |
//        |                      |
//        | ------- 栈底 -------- |

type controlStack struct {
	frames []*controlFrame
}

type controlFrame struct {
	// 创建当前帧的指令
	// 对于函数调用，创建帧的指令是 call
	// 对于流程控制所产生的帧，创建的指令有 block/loop 等
	opcode byte

	// 函数签名、以及块类型
	bt binary.FuncType

	// 复制了一份当前过程的指令
	instructions []binary.Instruction

	// base pointer 一个栈帧的开始的开始地址，对于函数调用来说，它是第 0 个实参的地址
	bp int

	// program counter 程序计数器，即当前指令的地址 **在当前帧** 里的位置，
	// 初始值为 0
	pc int
}

func newControlFrame(opcode byte,
	bt binary.FuncType,
	instructions []binary.Instruction,
	bp int) *controlFrame {
	// pc 初始值为 0
	return &controlFrame{opcode, bt, instructions, bp, 0}
}

func (s *controlStack) pushControlFrame(f *controlFrame) {
	s.frames = append(s.frames, f)
}

func (s *controlStack) popControlFrame() *controlFrame {
	lastIdx := len(s.frames) - 1
	f := s.frames[lastIdx]
	s.frames = s.frames[:lastIdx]
	return f
}

// -------- 辅助方法

func (s *controlStack) controlDepth() int { // name: frameSize
	return len(s.frames)
}

// 获取栈（包括控制栈和调用栈在内）顶的帧
func (s *controlStack) topControlFrame() *controlFrame { // name: topFrame
	return s.frames[len(s.frames)-1]
}

// 获取最后的一个**调用栈**的帧，即排除控制栈的帧。
// 返回调用帧以及该帧距离栈顶的距离（比如，如果栈顶帧就是调用帧，则距离为 0）
func (s *controlStack) topCallFrame() (*controlFrame, int) { // name: topCallFrame
	for idx := len(s.frames) - 1; idx >= 0; idx-- {
		if f := s.frames[idx]; f.opcode == binary.Call {
			return f, len(s.frames) - 1 - idx
		}
	}

	return nil, -1
}
