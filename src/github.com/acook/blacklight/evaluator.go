package main

func eval(ops []operation) {

	meta := NewMetaStack()
	current := NewSystemStack()
	meta.Push(current)

	doEval(meta, ops)
}

func doEval(meta *MetaStack, ops []operation) {
	var current *Stack
	defer rescue(meta)

	for _, op := range ops {
		s := meta.Peek()
		current = s.(*Stack)

		switch op.(type) {
		case *metaOp:
			meta = op.Eval(meta).(*MetaStack)
		case operation:
			current = op.Eval(current).(*Stack)
		default:
			panic("urecognized operation:" + op.String())
		}

	}
}

func rescue(meta *MetaStack) {
	if err := recover(); err != nil {
		warn("evaluation error")
		warn("$meta:", meta.String())
		s := meta.Peek()
		warn("@current:", s.String())
		panic(err)
	}
}

func conEval(name string, stack *Stack, ops []operation) {
	threads.Add(1)
	go func(name string, ops []operation, stack *Stack) {
		defer threads.Done()

		new_meta := NewMetaStack()
		new_meta.Push(stack)

		doEval(new_meta, ops)
	}(name, ops, stack)
}
