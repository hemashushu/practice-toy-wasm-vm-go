(module
    (func $f0
        (i32.const 10)
        (i32.const 0)
        (i32.eqz)
    )

    (func $f1
        (i32.const 10)
        (i32.const 1)
        ;; 不为 0 则压入 0
        (i32.eqz)
    )

    (func $f2
        (i32.const 10)
        (i32.const 20)
        ;; 不为 0 则压入 0
        (i32.eqz)
    )
)
