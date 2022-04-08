package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	count := len(args)

	if count == 2 {
		fmt.Println("Toy WebAssembly VM")
		fmt.Printf("running %s...\n", args[1])
		// exec args[1]
	} else {
		fmt.Println(`Toy WebAssembly VM
Usage:
$ go run . path_to_bytecode_file start_function_name

e.g.
$ go run . examples/01-hello.wasm hello`)
	}
}
