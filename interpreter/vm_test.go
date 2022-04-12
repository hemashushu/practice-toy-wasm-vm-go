package interpreter

import (
	"os"
	"path/filepath"
	"testing"
	"wasmvm/assert"
	"wasmvm/binary"
)

func TestInstConst(t *testing.T) {
	slots0 := run("test-vm-const.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{123}, slots0)

	slots1 := run("test-vm-const.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{123, 456}, slots1)

	slots2 := run("test-vm-const.wasm", 2)
	assert.AssertSliceEqual(t, []uint64{123}, slots2)

	slots3 := run("test-vm-const.wasm", 3)
	assert.AssertTrue(t, len(slots3) == 0)
}

func TestInstSelect(t *testing.T) {
	slots0 := run("test-vm-select.wasm", 0)
	assert.AssertSliceEqual(t, []uint64{100, 123}, slots0)

	slots1 := run("test-vm-select.wasm", 1)
	assert.AssertSliceEqual(t, []uint64{100, 456}, slots1)

	slots2 := run("test-vm-select.wasm", 2)
	assert.AssertSliceEqual(t, []uint64{100, 123}, slots2)
}

func run(fileName string, func_idx uint32) []uint64 {
	currentDir, err := os.Getwd() // Getwd() 返回当前 package 的目录，比如 `/path/to/project/binary`
	if err != nil {
		panic(err)
	}

	testResourcesDir := filepath.Join(currentDir, "..", "test", "resources", "interpreter")
	wasmFilePath := filepath.Join(testResourcesDir, fileName)

	m := binary.DecodeFile(wasmFilePath)
	return ExecFunc(m, func_idx)
}
