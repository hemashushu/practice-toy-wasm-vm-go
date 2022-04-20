(module
    (type $ft0 (func (param i32)))
    (type $ft1 (func (param i32 i32) (result i32)))
    ;; (import "env" "print_char" (func $print_char (type $ft0)))
    ;; (import "env" "print_int" (func $print_int (type $ft0)))
    (import "env" "add_i32" (func $add_i32 (type $ft1)))

;;     (func $f3
;;         (i32.const 65)
;;         (call $print_char)
;;     )
;;
;;     (func $f4
;;         (i32.const 65)
;;         (call $print_int)
;;     )

    ;; 导入了 1 个 native function，所以内部函数的
    ;; 索引值从 1 开始。

    (func $f1 (result i32)
        (i32.const 11)
        (i32.const 22)
        (call $add_i32)
    )
)
