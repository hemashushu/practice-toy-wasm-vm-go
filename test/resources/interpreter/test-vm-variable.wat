(module
    ;; todo
	(global $g1 (mut i32) (i32.const 1))
	(global $g2 (mut i32) (i32.const 2))
	(func (param $a i32) (param $b i32)
        (local $la i32)
        (local $lb i32)
        (local i64 i64)
		(global.get $g1)
		(global.set $g2)
		(local.get $a)
		(local.set $b)
	)
)
