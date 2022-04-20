package executor

import (
	"wasmvm/binary"
	"wasmvm/instance"
	"wasmvm/interpreter"
	"wasmvm/native"
)

func NewModule(m binary.Module) instance.Module {
	return NewModules([]string{"user"}, []binary.Module{m})["user"]
}

func NewModules(names []string, ms []binary.Module) map[string]instance.Module {
	moduleMap := map[string]instance.Module{}
	moduleMap["env"] = native.NewEnvModule()

	for idx, name := range names {
		moduleMap[name] = interpreter.NewModule(ms[idx], moduleMap)
	}

	return moduleMap
}
