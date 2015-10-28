package main

func eval(ops []operation) (bool, string) {

	meta := NewStack()
	current := NewStack()
	meta.Push(current)

	for _, op := range ops {
		switch op.(type) {
		case metaOp:
			op.Eval(meta)
		case operation:
			op.Eval(current)
		default:
			warn(op.String())
			panic("wait what")
		}
	}

	return true, ""
}
