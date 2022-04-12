(module
  (func $f0
      (i32.const 123)
  )

  (func $f1
      (i32.const 123)
      (i32.const 456)
  )

  (func $f2
      (i32.const 123)
      (i32.const 456)
      ;; drop 一次，弹出了 456
      (drop)
  )

  (func $f3
      (i32.const 123)
      (i32.const 456)
      ;; drop 两次，清空了操作数栈
      (drop)
      (drop)
  )
)
