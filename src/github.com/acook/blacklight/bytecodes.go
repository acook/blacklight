package main

// bytes spec

// (length)   : length specified by attribute
// (bytetype) : uint8 in range 0xF1-0xFA to specify type
// (type)     : expects full type signature
// (length types) : expects length number of full type signatures

// first byte (uint8)

// 0x00-0xF0 : opword (currently auto-generated)
// 0xF1 : word - uint32
// 0xF2 : byte - uint8 (byte)
// 0xF3 : char - uint32 (rune)
// 0xF4 : number - int64
// 0xF5 : float - float64
// 0xF6 : stack
// 0xF7 : queue
// 0xF8 : vector - kind:(bytetype) length:uint64 data:(length types)
// 0xF9 : object - length:int64 [slot:uint32 value:(type)]
// 0xFA : tag - kind:uint8 metadata:uint32 msg:CV
// 0xFB-0xFE : FUTURE DATATYPES
// 0xFF : RESERVED EXTENDED FLAG

var cv_map = map[string]byte{
	"opword": 0xF0,
	"word":   0xF1,
	"byte":   0xF2,
	"char":   0xF3,
	"number": 0xF4,
	"float":  0xF5,
	"stack":  0xf6,
}

var op_map map[string]byte
var fn_map map[byte]func(m *Meta)
var lk_map map[byte]string
var total_ops uint8

func prepare_op_table() {

	var op_fn_map = map[string]func(*Meta){

		// meta
		"@": func(m *Meta) {
			m.Current().Push(m.Current())
		},
		"^": func(m *Meta) {
			c := m.Current()
			n := m.Items[len(m.Items)-2]
			c.Push(n)
		},
		"$": func(m *Meta) {
			m.Current().Push(m)
		},
		"$decap": func(m *Meta) {
			m.Decap()
		},
		"$drop": func(m *Meta) {
			m.Drop()
			if m.Depth() < 1 {
				m.Push(NewSystemStack())
			}
		},
		"$new": func(m *Meta) {
			s := NewSystemStack()
			s.Push(m.Current())
			m.Push(s)
		},
		"$swap": func(m *Meta) {
			m.Swap()
		},

		// stack
		//"pop": func(m *Meta) {
		//	m.Current().Pop()
		//},
		//"drop": func(m *Meta) {
		//	m.Current().Drop()
		//},
		"decap": func(m *Meta) {
			m.Current().Decap()
		},
		"depth": func(m *Meta) {
			m.Current().Push(NewNumber(m.Current().Depth()))
		},
		"drop": func(m *Meta) {
			m.Current().Drop()
		},
		"dup": func(m *Meta) {
			m.Current().Dup()
		},
		"over": func(m *Meta) {
			m.Current().Over()
		},
		"purge": func(m *Meta) {
			m.Current().Purge()
		},
		"rot": func(m *Meta) {
			m.Current().Rot()
		},
		"swap": func(m *Meta) {
			m.Current().Swap()
		},

		// concurrency
		"bkg": func(m *Meta) {
			ops := m.Current().Pop().(WordVector).Ops

			items := NewStack("bkg")
			items.Push(m.Current().Pop())

			conEval("bkg", items, ops)
		},
		"co": func(m *Meta) {
			filename := m.Current().Pop().String()

			in := NewQueue()
			out := NewQueue()

			stack := NewStack("co")
			stack.Push(in)
			stack.Push(out)

			code := loadFile(filename)
			tokens := parse(code)
			ops := lex(tokens)

			m.Current().Push(out)
			m.Current().Push(in)

			conEval("co", stack, ops)
		},
		"work": func(m *Meta) {
			ops := m.Current().Pop().(WordVector).Ops
			in := m.Current().Items[len(m.Current().Items)-1].(*Queue)
			out := m.Current().Items[len(m.Current().Items)-2].(*Queue)

			stack := NewStack("work")
			stack.Push(in)
			stack.Push(out)

			m.Current().Push(out)
			m.Current().Push(in)

			conEval("work", stack, ops)
		},
		"wait": func(m *Meta) {
			threads.Wait()
		},

		// debug
		"print": func(m *Meta) {
			print(m.Current().Pop().String())
		},
		"refl": func(m *Meta) {
			NOPE("refl")
		},
		"warn": func(m *Meta) {
			NOPE("warn")
		},

		// loading
		"do": func(m *Meta) {
			filename := m.Current().Pop().String()

			code := loadFile(filename)
			tokens := parse(code)
			ops := lex(tokens)

			doBC(m, ops)
		},
		"imp": func(m *Meta) {
			NOPE("imp")
		},

		// math
		"add": func(m *Meta) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 + n1))
		},
		"sub": func(m *Meta) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 - n1))
		},
		"div": func(m *Meta) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 / n1))
		},
		"mod": func(m *Meta) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 % n1))
		},
		"mul": func(m *Meta) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 * n1))
		},
		"n-to-c": func(m *Meta) {
			m.Current().Push(NewCharFromString(m.Current().Pop().String()))
		},
		"n-to-cv": func(m *Meta) {
			m.Current().Push(NewCharVector(m.Current().Pop().String()))
		},

		// file io
		"read": func(m *Meta) {
			source := m.Current().Pop()
			q := m.Current().Peek().(*Queue)
			io := ReadIO(source, q)
			m.Current().Push(io)
		},
		"write": func(m *Meta) {
			dest := m.Current().Pop()
			q := m.Current().Peek().(*Queue)
			io := WriteIO(dest, q)
			m.Current().Push(io)
		},

		// logic & loops
		"either": func(m *Meta) {
			comp := m.Current().Pop().(WordVector).Ops
			iffalse := m.Current().Pop().(WordVector).Ops
			iftrue := m.Current().Pop().(WordVector).Ops
			doBC(m, comp)
			if m.Current().Pop().(*Tag).Kind == "true" {
				doBC(m, iftrue)
			} else {
				doBC(m, iffalse)
			}
		},
		"eq": func(m *Meta) {
			i1 := m.Current().Pop()
			i2 := m.Current().Peek()
			m.Current().Push(blEq(i1, i2))
		},
		"if": func(m *Meta) {
			comp := m.Current().Pop().(WordVector).Ops
			actn := m.Current().Pop().(WordVector).Ops
			doBC(m, comp)
			if m.Current().Pop().(*Tag).Kind == "true" {
				doBC(m, actn)
			}
		},
		"is": func(m *Meta) {
			NOPE("is")
		},
		"not": func(m *Meta) {
			var t *Tag
			i := m.Current().Pop()

			switch i.(type) {
			case *Tag:
				if i.(*Tag).Kind == "nil" {
					t = NewTrue("not")
				} else {
					t = NewNil("not")
				}
			default:
				t = NewNil("not")
			}

			m.Current().Push(t)
		},
		"until": func(m *Meta) {
			comp := m.Current().Pop().(WordVector).Ops
			actn := m.Current().Pop().(WordVector).Ops
		Until:
			for {
				doBC(m, comp)
				if m.Current().Pop().(*Tag).Kind == "true" {
					break Until
				}
				doBC(m, actn)
			}
		},
		"while": func(m *Meta) {
			comp := m.Current().Pop().(WordVector).Ops
			actn := m.Current().Pop().(WordVector).Ops
		While:
			for {
				doBC(m, comp)
				if m.Current().Pop().(*Tag).Kind != "true" {
					break While
				}
				doBC(m, actn)
			}
		},
		"loop": func(m *Meta) {
			actn := m.Current().Pop().(WordVector).Ops
			for {
				doBC(m, actn)
			}
		},

		// objects
		"o-new": func(m *Meta) {
			m.Current().Push(NewObject())
		},
		"self": func(m *Meta) {
			m.Self()
		},
		"child": func(m *Meta) {
			o := m.Object()
			child := NewChildObject(o)
			m.Current().Push(child)
		},
		"fetch": func(m *Meta) {
			slot := m.Current().Pop().(Word)
			o := m.Object()
			m.Current().Push(o.Fetch(slot))
		},
		"get": func(m *Meta) {
			slot := m.Current().Pop().(Word)
			m.Object().Get(m, slot)
		},
		"set": func(m *Meta) {
			slot := m.Current().Pop().(Word)
			i := m.Current().Pop()
			m.Object().Set(slot, i)
		},

		// queues
		"q-new": func(m *Meta) {
			m.Current().Push(NewQueue())
		},
		"deq": func(m *Meta) {
			i := m.Current().Peek().(*Queue).Dequeue()
			m.Current().Push(i)
		},
		"enq": func(m *Meta) {
			i := m.Current().Pop()
			m.Current().Peek().(*Queue).Enqueue(i)
		},
		"proq": func(m *Meta) {
			wv := m.Current().Pop().(WordVector)
			q := m.Current().Pop().(*Queue)

		ProcQLoop:
			for {
				select {
				case item := <-q.Items:
					m.Current().Push(item)
					doBC(m, wv.Ops)
				default:
					break ProcQLoop
				}
			}
		},
		"q-to-s": func(m *Meta) {
			NOPE("q-to-s")
		},
		"q-to-v": func(m *Meta) {
			q := m.Current().Pop().(*Queue)
			items := []datatypes{}

		QtoV:
			for {
				select {
				case i := <-q.Items:
					items = append(items, i)
				default:
					break QtoV
				}
			}

			m.Current().Push(NewVector(items))
		},
		"q-to-cv": func(m *Meta) {
			q := m.Current().Pop().(*Queue)
			str := ""

		QtoCV:
			for {
				i := <-q.Items
				if blEq(i, NewNil("q_to_cv")).Bool() {
					break QtoCV
				} else {
					str = str + i.(Char).CVString()
				}
			}

			v := NewCharVector(str)
			m.Current().Push(v)
		},
		"unq": func(m *Meta) {
			NOPE("unq")
		},

		// stacks
		"s-new": func(m *Meta) {
			m.Current().Push(NewStack("user"))
		},
		"pop": func(m *Meta) {
			m.Current().Push(m.Current().Peek().(stack).Pop())
		},
		"push": func(m *Meta) {
			m.Current().Peek().(stack).Push(m.Current().Pop())
		},
		"size": func(m *Meta) {
			m.Current().Push(NewNumber(m.Current().Peek().(stack).Depth()))
		},
		"tail": func(m *Meta) {
			m.Current().Peek().(stack).Drop()
		},

		// vectors
		"v-new": func(m *Meta) {
			m.Current().Push(NewVector([]datatypes{}))
		},
		"app": func(m *Meta) {
			i := m.Current().Pop()
			v := m.Current().Pop().(vector)
			m.Current().Push(v.App(i))
		},
		"ato": func(m *Meta) {
			c := m.Current()
			n := c.Pop().(N)
			v := m.Current().Peek().(vector)
			i := v.Ato(int(n.Value().(int64)))
			m.Current().Push(i)
		},
		"cat": func(m *Meta) {
			i1 := m.Current().Pop().(vector)
			i2 := m.Current().Pop().(vector)
			result := i2.Cat(i1)
			m.Current().Push(result)
		},
		"del": func(m *Meta) {
			NOPE("del")
		},
		"emt": func(m *Meta) {
			NOPE("emt")
		},
		"eval": func(m *Meta) {},
		"len": func(m *Meta) {
			v := m.Current().Peek().(vector)
			m.Current().Push(NewNumber(v.Len()))
		},
		"pick": func(m *Meta) {
			NOPE("pick")
		},
		"rmo": func(m *Meta) {
			n := m.Current().Pop().(*Number).Value().(int)
			v := m.Current().Pop().(vector)
			nv := v.Rmo(n)
			m.Current().Push(nv)
		},
		"v-to-s": func(m *Meta) {
			NOPE("v-to-s")
		},
		"v-to-q": func(m *Meta) {
			NOPE("v-to-q")
		},

		// tags
		"t-to-cv": func(m *Meta) {
			NOPE("t-to-cv")
		},
		"true": func(m *Meta) {
			m.Current().Push(NewTrue("true"))
		},
		"nil": func(m *Meta) {
			m.Current().Push(NewNil("nil"))
		},

		// chars
		"c-to-cv": func(m *Meta) {
			m.Current().Push(NewCharVector(m.Current().Pop().(Char).C_to_CV()))
		},
		"c-to-n": func(m *Meta) {
			m.Current().Push(NewNumber(m.Current().Pop().(Char).C_to_N()))
		},
	}

	op_map = make(map[string]byte)
	fn_map = make(map[byte]func(*Meta))
	lk_map = make(map[byte]string)

	var i byte = 0
	for k, v := range op_fn_map {
		op_map[k] = i
		fn_map[i] = v
		lk_map[i] = k
		i++
	}
	total_ops = uint8(i)

}

func NOPE(str string) {
	print(" -- UNIMPLEMENTED op: " + str + "\n")
}

func doBC(meta *Meta, ops []operation) {
	NOPE("can't call or eval shit yet")
}
