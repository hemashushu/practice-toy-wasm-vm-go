package executor

import (
	"os"
	"path/filepath"
	"testing"
	"wasmvm/assert"
	"wasmvm/binary"
	"wasmvm/instance"
)

func TestNativeFunction(t *testing.T) {
	// 测试调用本地函数（native function）
	assert.AssertListEqual(t,
		wrapList([]int32{33}), testFunc("test-executor-native-function.wasm", "test_add", nil))
}

func TestMultipleModules(t *testing.T) {
	assert.AssertListEqual(t,
		wrapList([]int32{77}),
		testFuncWithMultipleModules([]string{"lib", "app"}, "app", "test_add", nil))

	assert.AssertListEqual(t,
		wrapList([]int32{33}),
		testFuncWithMultipleModules([]string{"lib", "app"}, "app", "test_sub", nil))
}

func testFunc(fileName string, funcName string, args []instance.WasmVal) []instance.WasmVal {
	m := readModule(fileName)
	mod := NewModule(m)
	return mod.EvalFunc(funcName, args...)
}

func testFuncWithMultipleModules(
	moduleNames []string,
	targetModuleName string,
	targetFuncName string,
	args []instance.WasmVal) []instance.WasmVal {

	ms := []binary.Module{}
	for _, n := range moduleNames {
		fileName := "test-module-" + n + ".wasm"
		m := readModule(fileName)
		ms = append(ms, m)
	}

	mods := NewModules(moduleNames, ms)
	mod := mods[targetModuleName]
	return mod.EvalFunc(targetFuncName, args...)
}

func readModule(fileName string) binary.Module {
	currentDir, err := os.Getwd() // Getwd() 返回当前 package 的目录，比如 `/path/to/project/binary`
	if err != nil {
		panic(err)
	}

	testResourcesDir := filepath.Join(currentDir, "..", "test", "resources", "executor")
	wasmFilePath := filepath.Join(testResourcesDir, fileName)

	return binary.DecodeFile(wasmFilePath)
}

func wrapList[T comparable](items []T) []instance.WasmVal {
	rs := []instance.WasmVal{}
	for _, i := range items {
		rs = append(rs, i)
	}
	return rs
}
