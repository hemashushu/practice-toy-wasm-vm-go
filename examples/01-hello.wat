(module
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "hello" (func $hello))
 (func $hello (; 0 ;) (result i32)
  (i32.const 100)
 )
)
