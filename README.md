# (Practice) Toy WebAssembly VM - Go

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [(Practice) Toy WebAssembly VM - Go](#practice-toy-webassembly-vm-go)
  - [使用方法](#使用方法)
    - [测试](#测试)
    - [编译](#编译)
    - [运行指定的脚本](#运行指定的脚本)
  - [附录](#附录)
    - [工具之 wasm-tools](#工具之-wasm-tools)
      - [文本和二进制相互转换](#文本和二进制相互转换)
      - [查看二进制信息](#查看二进制信息)
      - [单单查看段信息](#单单查看段信息)
    - [工具之 wabt](#工具之-wabt)
      - [文本格式和二进制格式转换](#文本格式和二进制格式转换)
      - [查看二进制文件内容](#查看二进制文件内容)
      - [运行字节码](#运行字节码)
    - [从 Rust 编译到 Wasm](#从-rust-编译到-wasm)
      - [编译单独一个 Rust 源码文件](#编译单独一个-rust-源码文件)
      - [编译一个 cargo 项目](#编译一个-cargo-项目)
    - [从 C 编译到 wasm](#从-c-编译到-wasm)
    - [使用 wasm-pack 和 wasm-bindgen](#使用-wasm-pack-和-wasm-bindgen)

<!-- /code_chunk_output -->

练习单纯使用 Go lang 编写简单的 _WebAssembly VM_。

> 注：本项目是阅读和学习《WebAssembly 原理与核心技术》时的随手练习，并无实际用途。程序的原理、讲解和代码的原始出处请参阅书本。

## 使用方法

### 测试

`$ go test ./...`

### 编译

`$ go build -o wasmvm`

### 运行指定的脚本

`$ ./wasmvm path_to_bytecode_file start_function_name`

或者：

`$ go run . path_to_bytecode_file start_function_name`

示例：

`$ go run . examples/03-simple.wasm main`

如无意外，应该能看到输出 `3`。

## 附录

### 工具之 wasm-tools

https://github.com/bytecodealliance/wasm-tools

安装：

`$ cargo install wasm-tools`

#### 文本和二进制相互转换

文本转二进制：

`$ wasm-tools parse 01-hello.wat -o 01-hello.wasm`

二进制转文本：

`wasm-tools print 01-hello.wasm`

#### 查看二进制信息

显示二进制和信息的对照文本，相当于反编译

`$ wasm-tools dump 01-hello.wasm`

#### 单单查看段信息

`$ wasm-tools objdump 01-hello.wasm`

### 工具之 wabt

wabt

https://github.com/WebAssembly/wabt

#### 文本格式和二进制格式转换

命令示例：

- `$ wat2wasm hello.wat`
- `$ wasm2wat hello.wasm`

#### 查看二进制文件内容

`$ wasm-objdump -h hello.wasm`

可选参数：

- -h, --headers
  Print headers
- -j, --section=SECTION
  Select just one section
- -s, --full-contents
  Print raw section contents
- -d, --disassemble
  Disassemble function bodies

#### 运行字节码

`$ wasm-interp test.wasm --run-all-exports`

### 从 Rust 编译到 Wasm

先添加 target `wasm32-unknown-unknown`

`$ rustup target add wasm32-unknown-unknown`

参考：

https://www.hellorust.com/setup/wasm-target/

#### 编译单独一个 Rust 源码文件

`$ rustc --target wasm32-unknown-unknown -C lto -O --crate-type=cdylib 02-rust.rs -o 02-rust.wasm`

#### 编译一个 cargo 项目

确保 `Cargo.toml` 的内容如下：

```toml
[lib]
path = "src/lib.rs"
crate-type = ["cdylib"]
```

然后开始编译：

`$ cargo build --target wasm32-unknown-unknown --release`

### 从 C 编译到 wasm

`$ clang --target=wasm32 -c --no-standard-libraries 01-hello.c -o 01-hello.wasm`

参考：

- https://depth-first.com/articles/2019/10/16/compiling-c-to-webassembly-and-running-it-without-emscripten/
- https://dassur.ma/things/c-to-webassembly/
- https://00f.net/2019/04/07/compiling-to-webassembly-with-llvm-and-clang/

### 使用 wasm-pack 和 wasm-bindgen

构建跟 JavaScript 互动的程序

参考：

https://developer.mozilla.org/en-US/docs/WebAssembly/Rust_to_wasm
