(module
    ;; 测试 eq 和 ne

    (func $f0
        (i32.const 10)
        (i32.const 20)
        (i32.const 30)
        (i32.eq)        ;; false/0
    )

    (func $f1
        (i32.const 10)
        (i32.const 20)
        (i32.const 30)
        (i32.ne)        ;; true/1
    )

    (func $f2
        (i32.const 10)
        (i32.const 20)
        (i32.const 20)
        (i32.eq)        ;; true/1
    )

    (func $f3
        (i32.const 10)
        (i32.const 20)
        (i32.const 20)
        (i32.ne)        ;; false/0
    )

    ;; 测试 lt, gt

    (func $f4
        (i32.const 10)
        (i32.const 20)  ;; 后弹出，LHS
        (i32.const -30) ;; 先弹出，RHS
        (i32.lt_s)      ;; false/0
    )

    (func $f5
        (i32.const 10)
        (i32.const 20)
        (i32.const -30)
        ;; 将 -30 作为 unsigned 整数时，比 20 大
        (i32.lt_u)      ;; true/1
    )

    (func $f6
        (i32.const 10)
        (i32.const 20)
        (i32.const -30)
        (i32.gt_s)      ;; true/1
    )

    (func $f7
        (i32.const 10)
        (i32.const 20)
        (i32.const -30)
        (i32.gt_u)      ;; false/0
    )

    ;; 测试 le, ge

    (func $f8
        (i32.const 10)
        (i32.const 20)  ;; 后弹出，LHS
        (i32.const -30) ;; 先弹出，RHS
        (i32.le_s)      ;; false/0
    )

    (func $f9
        (i32.const 10)
        (i32.const 20)
        (i32.const -30)
        ;; 将 -30 作为 unsigned 整数时，比 20 大
        (i32.le_u)      ;; true/1
    )

    (func $f10
        (i32.const 10)
        (i32.const 20)
        (i32.const -30)
        (i32.ge_s)      ;; true/1
    )

    (func $f11
        (i32.const 10)
        (i32.const 20)
        (i32.const -30)
        (i32.ge_u)      ;; false/0
    )

    ;; 测试 le, ge（当两个操作数相等的情况）

    (func $f12
        (i32.const 10)
        (i32.const 20)  ;; 后弹出，LHS
        (i32.const 20)  ;; 先弹出，RHS
        (i32.le_s)      ;; true/1
    )

    (func $f13
        (i32.const 10)
        (i32.const 20)
        (i32.const 20)
        (i32.le_u)      ;; true/1
    )

    (func $f14
        (i32.const 10)
        (i32.const 20)
        (i32.const 20)
        (i32.ge_s)      ;; true/1
    )

    (func $f15
        (i32.const 10)
        (i32.const 20)
        (i32.const 20)
        (i32.ge_u)      ;; true/1
    )

    ;; 测试浮点数

    (func $f16
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 3.3)
        (f32.eq)        ;; false/0
    )

    (func $f17
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 3.3)
        (f32.ne)        ;; true/1
    )

    (func $f18
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 2.2)
        (f32.eq)        ;; true/1
    )

    (func $f19
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 2.2)
        (f32.ne)        ;; false/0
    )

    (func $f20
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 3.3)
        (f32.lt)        ;; true/1
    )

    (func $f21
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 3.3)
        (f32.gt)        ;; false/0
    )

    (func $f22
        (i32.const 11)
        (f32.const 2.2)
        (f32.const -3.3)
        (f32.lt)        ;; false/0
    )

    (func $f23
        (i32.const 11)
        (f32.const 2.2)
        (f32.const -3.3)
        (f32.gt)        ;; true/1
    )

    ;; 测试 le, ge

    (func $f24
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 3.3)
        (f32.le)        ;; true/1
    )

    (func $f25
        (i32.const 11)
        (f32.const 2.2)
        (f32.const -3.3)
        (f32.ge)        ;; false/0
    )

    ;; 测试 le, ge（当两个操作数相等的情况）

    (func $f26
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 2.2)
        (f32.le)        ;; true/1
    )

    (func $f27
        (i32.const 11)
        (f32.const 2.2)
        (f32.const 2.2)
        (f32.ge)        ;; true/1
    )

)
