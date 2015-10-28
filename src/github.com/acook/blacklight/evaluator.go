package main

func eval(ops []operation) (bool, string) {

	meta := *NewStack()
	current := *NewStack()
	meta.Push(current)

	defer rescue(meta)

	for _, op := range ops {
		switch op.(type) {
		case metaOp:
			meta = op.Eval(meta)
		case operation:
			current = op.Eval(current)
		default:
			warn(op.String())
			panic("wait what")
		}
	}

	return true, ""
}

func rescue(meta Stack) {
	if err := recover(); err != nil {
		warn("evaluation error")
		warn("$meta:", meta.String())
		warn("@current:", meta.Peek().String())
		panic(err)
	}
}
