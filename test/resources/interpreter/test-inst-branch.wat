(module
    (func
        (block
            (block
                (block
                    (br 1)
                    (br_if 2 (i32.const 100))
                    (br_table 0 1 2 3) ;; 3 是默认标签
                    (return)
                )
            )
        )
    )
)

;; 0x0015 | 19          | size of function
;; 0x0016 | 00          | 0 local blocks
;; 0x0017 | 02 40       | Block { ty: Empty }
;; 0x0019 | 02 40       |   Block { ty: Empty }
;; 0x001b | 02 40       |     Block { ty: Empty }
;; 0x001d | 0c 01       |       Br { relative_depth: 1 }
;; 0x001f | 41 e4 00    |       I32Const { value: 100 }
;; 0x0022 | 0d 02       |       BrIf { relative_depth: 2 }
;; 0x0024 | 0e 03 00 01 |       BrTable { table: BrTable { count: 3, default: 3, targets: [0, 1, 2] } }
;;        | 02 03
;; 0x002a | 0f          |       Return
;; 0x002b | 0b          |       End
;; 0x002c | 0b          |     End
;; 0x002d | 0b          |   End
;; 0x002e | 0b          | End