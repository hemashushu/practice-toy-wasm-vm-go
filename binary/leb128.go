package binary

import "errors"

// leb128 整数编码
//
// 一种变长的整数编码格式，int32 编码后有 1~5 字节长
// int64 编码后有 1~10 字节长
// https://en.wikipedia.org/wiki/LEB128
//
// 原理：
// 每个字节如果最高位（第 7 位，从 0 开始数）为 1，则表示还有后续的
// 内容。比如
//
// byte 0     byte 1     byte 2
// -----------------  ---------  ---------
// 1000 0001, 1000 0010, 0000 0000
// ^          ^          ^
// |--有后续   |--有后续   |--无后续
//
// 每个字节的低 7 位拼接在一起就是最终的整数：
//
// 000 0000 000 0010 000 0001
// ---------------- -------- --------
// byte 2   byte 1   byte 0
//
// 对于有符号整数，拼接后的最高位是符号位，比如对于一个 int16 数：
//
// byte 0     byte 1
// -----------------  ---------
// 1000 0001, 0111 0010
// ^          ^^
// |          ||--- 符号位
// |--有后续   |---- 无后续
//
// 拼接后的整数：
//
// 111 0010 000 0001
// ---------------- --------
// byte 1   byte 0
//
// 最高位是 1，所以所有高位都需要补上 1，最后得整数（int16）
//
// 1111 1001 0000 0001
// ^^---补上的 1

// 解码 uint32 或者 uint64
//
// 返回：
// 1. 解码后的整数，如需解码 uint32，则强行将返回值转为 uint32
// 2. 实际消耗掉的字节数
func decodeVarUint(data []byte, bitWidth int) (uint64, int) {
	result := uint64(0)
	for i, b := range data {
		if i == bitWidth/7 {
			// 到达最后一个字节
			if b&0b1000_0000 != 0 {
				// 最后一个字节的索引 7 比特应该为 0
				panic(errors.New("data too long"))
			}

			if b>>(bitWidth-i*7) != 0 {
				// 超出 bitWidth 的部分应该为 0，否则就溢出了
				panic(errors.New("int overflow"))
			}
		}

		result |= (uint64(b) & 0b0111_1111) << (i * 7)
		if b&0b1000_0000 == 0 {
			return result, i + 1
		}
	}

	panic(errors.New("unexpected the end of leb128"))
}

func decodeVarInt(data []byte, bitWidth int) (int64, int) {
	result := int64(0)
	for i, b := range data {
		if i == bitWidth/7 {
			// 到达最后一个字节
			if b&0b10000000 != 0 {
				// 最后一个字节的索引 7 比特应该为 0
				panic(errors.New("data too long"))
			}

			if b&0b0100_0000 == 0 && b>>(bitWidth-i*7-1) != 0 {
				// 当前为正数，则超出 bitWidth 的部分应该为 0，否则就溢出了
				panic(errors.New("int overflow"))
			}

			if b&0b0100_0000 != 0 && int8(b|0b1000_0000)>>(bitWidth-i*7-1) != -1 {
				// 当前为负数，则超出 bitWidth 的部分应该全为 1，否则不合法
				panic(errors.New("invalid negative"))
			}
		}

		result |= (int64(b) & 0b0111_1111) << (i * 7)
		if b&0b1000_0000 == 0 {
			if (i*7 < bitWidth) && (b&0b0100_0000 != 0) {
				// 如果索引 6 比特为 1，表示负数，需要把高位都补上 1
				result = result | (-1 << ((i + 1) * 7))
			}
			return result, i + 1
		}
	}

	panic(errors.New("unexpected the end of leb128"))
}
