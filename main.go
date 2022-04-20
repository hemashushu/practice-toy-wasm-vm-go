package main

import (
	"fmt"
	"os"
	"path/filepath"
	"wasmvm/binary"
	"wasmvm/executor"
)

func main() {
	args := os.Args
	count := len(args)

	if count == 3 {
		fmt.Println("Toy WebAssembly VM")
		fmt.Printf("running %s, func: %s ...\n", args[1], args[2])
		// exec args[1]
		exec(args[1], args[2])
	} else {
		fmt.Println(`Toy WebAssembly VM
Usage:
$ go run . path_to_bytecode_file start_function_name

e.g.
$ go run . examples/03-simple.wasm main`)
	}
}

func exec(fileName string, funcName string) {
	currentDir, err := os.Getwd() // Getwd() 返回当前 package 的目录，比如 `/path/to/project`
	if err != nil {
		panic(err)
	}

	wasmFilePath := filepath.Join(currentDir, fileName)

	m := binary.DecodeFile(wasmFilePath)
	mod := executor.NewModule(m)
	r := mod.EvalFunc(funcName)
	fmt.Printf("%v\n", r)
}
