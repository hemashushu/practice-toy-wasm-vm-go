(module
  (type (;0;) (func (param i32) (result i32)))
  (func $add_one (type 0) (param i32) (result i32)
    local.get 0
    i32.const 1
    i32.add
  )
  (memory (;0;) 16)
  (global $__stack_pointer (mut i32) i32.const 1048576)
  (global (;1;) i32 i32.const 1048576)
  (global (;2;) i32 i32.const 1048576)
  (export "memory" (memory 0))
  (export "add_one" (func $add_one))
  (export "__data_end" (global 1))
  (export "__heap_base" (global 2))
)
