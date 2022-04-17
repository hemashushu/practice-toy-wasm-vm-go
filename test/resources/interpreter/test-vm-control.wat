(module
    (func $0 (result i32)
        (i32.const 1)
        (return)
        (i32.const 10)
    )

    (func $1 (result i32)
        (i32.const 1)
        (block
            (i32.const 2)
            (return)
        )
        (i32.const 10)
    )

    (func $2 (result i32)
        (i32.const 1)
        (block
            (i32.const 2)
            (block
                (i32.const 3)
                (return)
            )
        )
        (i32.const 10)
    )

)
