package executor

import (
	"wasmvm/binary"
	"wasmvm/instance"
	"wasmvm/interpreter"
	"wasmvm/native"
)

func newModule(m binary.Module) instance.Module {
	return newModules([]string{"user"}, []binary.Module{m})["user"]
}

func newModules(names []string, ms []binary.Module) map[string]instance.Module {
	moduleMap := map[string]instance.Module{}
	moduleMap["env"] = native.NewEnvModule()

	for idx, name := range names {
		moduleMap[name] = interpreter.NewModule(ms[idx], moduleMap)
	}

	return moduleMap
}
