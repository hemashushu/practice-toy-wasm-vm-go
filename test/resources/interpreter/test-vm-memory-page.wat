(module
    (memory 2 10)
    (func $f0 (result i32 i32)
        (i32.const 10)
        (memory.size)
    )
    (func $f1 (result i32 i32 i32 i32)
        (i32.const 10)
        (i32.const 2)
        (memory.grow)
        (i32.const 3)
        (memory.grow)
        (memory.size)
    )
    ;; 注：这里没有测试页面增加后，原有的数据是否仍然存在
)
