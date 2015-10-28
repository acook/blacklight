package main

func eval(ops []operation) (bool, string) {

	meta := NewMetaStack()
	current := NewSystemStack()
	meta.Push(current)

	defer rescue(meta)

	for _, op := range ops {
		s := *meta.Peek()
		current = s.(*SystemStack)

		switch op.(type) {
		case *metaOp:
			meta = op.Eval(meta).(*MetaStack)
		case operation:
			current = op.Eval(current).(*SystemStack)
		default:
			warn(op.String())
			panic("wait what")
		}

	}

	return true, ""
}

func rescue(meta stack) {
	if err := recover(); err != nil {
		warn("evaluation error")
		warn("$meta:", meta.String())
		s := *meta.Peek()
		warn("@current:", s.String())
		panic(err)
	}
}
