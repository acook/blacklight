package main

// bytes spec

// (chain)    : max value (or sign) means additional byte continues value, following byte may be 0x00 if value is actually 0xFF
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
// 0xF8 : vector - kind:(bytetype) length:int8(chain) data:(length types)
// 0xF9 : object - length:int8(chain) [slot:uint32 value:(type)]
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
var fn_map map[byte]func(m *MetaStack)

func prepare_op_table() {

	var op_fn_map = map[string]func(*MetaStack){

		// meta
		"@": func(m *MetaStack) {
			m.Current().Push(m.Current())
		},
		"^": func(m *MetaStack) {
			c := m.Current()
			n := m.Items[len(m.Items)-2]
			c.Push(n)
		},
		"$": func(m *MetaStack) {
			m.Current().Push(m)
		},
		"$decap": func(m *MetaStack) {
			m.Decap()
		},
		"$drop": func(m *MetaStack) {
			m.Drop()
			if m.Depth() < 1 {
				m.Push(NewSystemStack())
			}
		},
		"$new": func(m *MetaStack) {
			s := NewSystemStack()
			s.Push(m.Current())
			m.Push(s)
		},
		"$swap": func(m *MetaStack) {
			m.Swap()
		},

		// stack
		//"pop": func(m *MetaStack) {
		//	m.Current().Pop()
		//},
		//"drop": func(m *MetaStack) {
		//	m.Current().Drop()
		//},
		"decap": func(m *MetaStack) {
			m.Current().Decap()
		},
		"depth": func(m *MetaStack) {
			m.Current().Push(NewNumber(m.Current().Depth()))
		},
		"drop": func(m *MetaStack) {
			m.Current().Drop()
		},
		"dup": func(m *MetaStack) {
			m.Current().Dup()
		},
		"over": func(m *MetaStack) {
			m.Current().Over()
		},
		"purge": func(m *MetaStack) {
			m.Current().Purge()
		},
		"rot": func(m *MetaStack) {
			m.Current().Rot()
		},
		"swap": func(m *MetaStack) {
			m.Current().Swap()
		},

		// concurrency
		"bkg": func(m *MetaStack) {
			ops := m.Current().Pop().(WordVector).Ops

			items := NewStack("bkg")
			items.Push(m.Current().Pop())

			conEval("bkg", items, ops)
		},
		"co": func(m *MetaStack) {
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
		"work": func(m *MetaStack) {
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
		"wait": func(m *MetaStack) {
			threads.Wait()
		},

		// debug
		"print": func(m *MetaStack) {
			print(m.Current().Pop().String())
		},
		"refl": func(m *MetaStack) {
			NOPE("refl")
		},
		"warn": func(m *MetaStack) {
			NOPE("warn")
		},

		// loading
		"do": func(m *MetaStack) {
			filename := m.Current().Pop().String()

			code := loadFile(filename)
			tokens := parse(code)
			ops := lex(tokens)

			doEval(m, ops)
		},
		"imp": func(m *MetaStack) {
			NOPE("imp")
		},

		// math
		"add": func(m *MetaStack) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 + n1))
		},
		"sub": func(m *MetaStack) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 - n1))
		},
		"div": func(m *MetaStack) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 / n1))
		},
		"mod": func(m *MetaStack) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 % n1))
		},
		"mul": func(m *MetaStack) {
			n1 := m.Current().Pop().Value().(int)
			n2 := m.Current().Pop().Value().(int)

			m.Current().Push(NewNumber(n2 * n1))
		},
		"n-to-c": func(m *MetaStack) {
			m.Current().Push(NewCharFromString(m.Current().Pop().String()))
		},
		"n-to-cv": func(m *MetaStack) {
			m.Current().Push(NewCharVector(m.Current().Pop().String()))
		},

		// file io
		"read": func(m *MetaStack) {
			source := m.Current().Pop()
			q := m.Current().Peek().(*Queue)
			io := ReadIO(source, q)
			m.Current().Push(io)
		},
		"write": func(m *MetaStack) {
			dest := m.Current().Pop()
			q := m.Current().Peek().(*Queue)
			io := WriteIO(dest, q)
			m.Current().Push(io)
		},

		// logic & loops
		"either": func(m *MetaStack) {
			comp := m.Current().Pop().(WordVector).Ops
			iffalse := m.Current().Pop().(WordVector).Ops
			iftrue := m.Current().Pop().(WordVector).Ops
			doEval(m, comp)
			if m.Current().Pop().(*Tag).Kind == "true" {
				doEval(m, iftrue)
			} else {
				doEval(m, iffalse)
			}
		},
		"eq": func(m *MetaStack) {
			i1 := m.Current().Pop()
			i2 := m.Current().Peek()
			m.Current().Push(blEq(i1, i2))
		},
		"if": func(m *MetaStack) {
			comp := m.Current().Pop().(WordVector).Ops
			actn := m.Current().Pop().(WordVector).Ops
			doEval(m, comp)
			if m.Current().Pop().(*Tag).Kind == "true" {
				doEval(m, actn)
			}
		},
		"is": func(m *MetaStack) {
			NOPE("is")
		},
		"not": func(m *MetaStack) {
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
		"until": func(m *MetaStack) {
			comp := m.Current().Pop().(WordVector).Ops
			actn := m.Current().Pop().(WordVector).Ops
		Until:
			for {
				doEval(m, comp)
				if m.Current().Pop().(*Tag).Kind == "true" {
					break Until
				}
				doEval(m, actn)
			}
		},
		"while": func(m *MetaStack) {
			comp := m.Current().Pop().(WordVector).Ops
			actn := m.Current().Pop().(WordVector).Ops
		While:
			for {
				doEval(m, comp)
				if m.Current().Pop().(*Tag).Kind != "true" {
					break While
				}
				doEval(m, actn)
			}
		},
		"loop": func(m *MetaStack) {
			actn := m.Current().Pop().(WordVector).Ops
			for {
				doEval(m, actn)
			}
		},

		// objects
		"o-new": func(m *MetaStack) {
			m.Current().Push(NewObject())
		},
		"self": func(m *MetaStack) {
			m.Self()
		},
		"child": func(m *MetaStack) {
			o := m.Object()
			child := NewChildObject(o)
			m.Current().Push(child)
		},
		"fetch": func(m *MetaStack) {
			slot := m.Current().Pop().(Word)
			o := m.Object()
			m.Current().Push(o.Fetch(slot))
		},
		"get": func(m *MetaStack) {
			slot := m.Current().Pop().(Word)
			m.Object().Get(m, slot)
		},
		"set": func(m *MetaStack) {
			slot := m.Current().Pop().(Word)
			i := m.Current().Pop()
			m.Object().Set(slot, i)
		},

		// queues
		"q-new": func(m *MetaStack) {
			m.Current().Push(NewQueue())
		},
		"deq": func(m *MetaStack) {
			i := m.Current().Peek().(*Queue).Dequeue()
			m.Current().Push(i)
		},
		"enq": func(m *MetaStack) {
			i := m.Current().Pop()
			m.Current().Peek().(*Queue).Enqueue(i)
		},
		"proq": func(m *MetaStack) {
			wv := m.Current().Pop().(WordVector)
			q := m.Current().Pop().(*Queue)

		ProcQLoop:
			for {
				select {
				case item := <-q.Items:
					m.Current().Push(item)
					doEval(m, wv.Ops)
				default:
					break ProcQLoop
				}
			}
		},
		"q-to-s": func(m *MetaStack) {
			NOPE("q-to-s")
		},
		"q-to-v": func(m *MetaStack) {
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
		"q-to-cv": func(m *MetaStack) {
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
		"unq": func(m *MetaStack) {
			NOPE("unq")
		},

		// stacks
		"s-new": func(m *MetaStack) {
			m.Current().Push(NewStack("user"))
		},
		"pop": func(m *MetaStack) {
			m.Current().Push(m.Current().Peek().(stack).Pop())
		},
		"push": func(m *MetaStack) {
			m.Current().Peek().(stack).Push(m.Current().Pop())
		},
		"size": func(m *MetaStack) {
			m.Current().Push(NewNumber(m.Current().Peek().(stack).Depth()))
		},
		"tail": func(m *MetaStack) {
			m.Current().Peek().(stack).Drop()
		},

		// vectors
		"v-new": func(m *MetaStack) {
			m.Current().Push(NewVector([]datatypes{}))
		},
		"app": func(m *MetaStack) {
			i := m.Current().Pop()
			v := m.Current().Pop().(vector)
			m.Current().Push(v.App(i))
		},
		"ato": func(m *MetaStack) {
			n := m.Current().Pop().(*Number)
			v := m.Current().Peek().(vector)
			i := v.Ato(n.Value().(int))
			m.Current().Push(i)
		},
		"cat": func(m *MetaStack) {
			i1 := m.Current().Pop().(vector)
			i2 := m.Current().Pop().(vector)
			result := i2.Cat(i1)
			m.Current().Push(result)
		},
		"del": func(m *MetaStack) {
			NOPE("del")
		},
		"emt": func(m *MetaStack) {
			NOPE("emt")
		},
		"eval": func(m *MetaStack) {},
		"len": func(m *MetaStack) {
			v := m.Current().Peek().(vector)
			m.Current().Push(NewNumber(v.Len()))
		},
		"pick": func(m *MetaStack) {
			NOPE("pick")
		},
		"rmo": func(m *MetaStack) {
			n := m.Current().Pop().(*Number).Value().(int)
			v := m.Current().Pop().(vector)
			nv := v.Rmo(n)
			m.Current().Push(nv)
		},
		"v-to-s": func(m *MetaStack) {
			NOPE("v-to-s")
		},
		"v-to-q": func(m *MetaStack) {
			NOPE("v-to-q")
		},

		// tags
		"t-to-cv": func(m *MetaStack) {
			NOPE("t-to-cv")
		},
		"true": func(m *MetaStack) {
			m.Current().Push(NewTrue("true"))
		},
		"nil": func(m *MetaStack) {
			m.Current().Push(NewNil("nil"))
		},

		// chars
		"c-to-cv": func(m *MetaStack) {
			m.Current().Push(NewCharVector(m.Current().Pop().(Char).C_to_CV()))
		},
		"c-to-n": func(m *MetaStack) {
			m.Current().Push(NewNumber(m.Current().Pop().(Char).C_to_N()))
		},
	}

	op_map = make(map[string]byte)
	fn_map = make(map[byte]func(*MetaStack))

	var i byte = 0
	for k, v := range op_fn_map {
		op_map[k] = i
		fn_map[i] = v
		i++
	}

}

func NOPE(str string) {
	print(" -- UNIMPLEMENTED op: " + str)
}
