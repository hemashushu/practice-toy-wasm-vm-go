(module
    (type $ft0 (func (param i32 i32) (result i32)))
    (export "add" (func $add))
    (export "sub" (func $sub))

    (func $add (type $ft0)
        (local.get 0)
        (local.get 1)
        (i32.add)
    )

    (func $sub (type $ft0)
        (local.get 0)
        (local.get 1)
        (i32.sub)
    )
)
