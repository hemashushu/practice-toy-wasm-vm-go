(module
    (func $f0
        (i32.const 0)
        (i32.load8_u offset=0)
    )

    ;; 测试只改变 offset 立即数
    (func $f1
        (i32.const 0)
        (i32.load8_u offset=0)
        (i32.const 0)
        (i32.load8_u offset=1)
        (i32.const 0)
        (i32.load8_u offset=2)
        (i32.const 0)
        (i32.load8_u offset=3)
    )

    ;; 测试只改变地址
    (func $f2
        (i32.const 0)
        (i32.load8_u offset=0)
        (i32.const 1)
        (i32.load8_u offset=0)
        (i32.const 2)
        (i32.load8_u offset=0)
        (i32.const 3)
        (i32.load8_u offset=0)
    )

    ;; 测试以符号数来加载
    (func $f3
        (i32.const 0)
        (i32.load8_u offset=0) ;; 17
        (i32.const 0)
        (i32.load8_s offset=0) ;; 17
        (i32.const 1)
        (i32.load8_u offset=0) ;; 241
        (i32.const 1)
        (i32.load8_s offset=0) ;; -15
    )

    ;; 测试加载 16 位, 32 位数
    (func $4
        (i32.const 2)
        (i32.load16_u) ;; 0x6655
        (i32.const 2)
        (i32.load16_s) ;; 0x6655
        (i32.const 4)
        (i32.load16_u) ;; 0x9080
        (i32.const 4)
        (i32.load16_s) ;; 0x9080

        (i32.const 6)
        (i32.load)     ;; 32 位
    )

    ;; 测试加载 64 位数
    (func $5
        (i32.const 6)
        (i64.load32_u) ;; 0x0706050403020100
        (i32.const 6)
        (i64.load32_s) ;;
        (i32.const 14)
        (i64.load32_u) ;; 0xf0e0d0c0b0a09080
        (i32.const 14)
        (i64.load32_s) ;;

        (i32.const 6)
        (i64.load)     ;; 64 位
        (i32.const 14)
        (i64.load)     ;; 64 位
    )
)
