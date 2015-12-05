package main

import (
	"encoding/binary"
	//"fmt"
)

// bytecode type functions
// signature: (vm *VMState) void

func NOPE(str string) {
	print(" -- UNIMPLEMENTED op: " + str + "\n")
}

func nullbyte(vm *VMstate) {
	print(" -- NULL BYTE ENCOUNTERED\n")
	vm.debug()
	print("\n")
	panic("vm: aw shit, something is terribly wrong, we encountered a null byte")
}

func opword(vm *VMstate) {
	fn_map[vm.b](vm.m)
}

func word(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	wd_uint := binary.BigEndian.Uint64(buf)

	vm.o = vm.o + 7
	vm.m.Current().Push(W(wd_uint))
}

func octet(vm *VMstate) {
	vm.o++
	c := vm.bc[vm.o]
	vm.m.Current().Push(C(c))
}

func bl_rune(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+4]

	r := Varint32(buf)

	vm.o = vm.o + 3
	vm.m.Current().Push(R(r))
}

func integer(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]

	n := Varint64(buf)

	vm.o = vm.o + 7
	vm.m.Current().Push(N(n))
}

func text(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	length := binary.BigEndian.Uint64(buf)
	vm.o = vm.o + 7

	vm.o++
	str_buf := vm.bc[vm.o : vm.o+length]

	vm.o = vm.o + (length - 1)
	vm.m.Current().Push(T(str_buf))
}

func block(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	length := binary.BigEndian.Uint64(buf)
	vm.o = vm.o + 7

	vm.o++
	blk_buf := vm.bc[vm.o : vm.o+length]

	vm.o = vm.o + (length - 1)
	vm.m.Current().Push(B(blk_buf))
}

func vector(vm *VMstate) {
	vm.m.NewStack("vector")
}

func endvector(vm *VMstate) {
	v := vm.m.Eject().S_to_V()
	vm.m.Current().Push(v)
}

// opword instructions
// signature: (m *Meta) void

// META OPS

func push_current(m *Meta) {
	m.Current().Push(m.Current())
}

func push_last(m *Meta) {
	c := m.Current()
	n := m.Items[len(m.Items)-2]
	c.Push(n)
}

func push_meta(m *Meta) {
	m.Current().Push(m)
}

func meta_decap(m *Meta) {
	m.Decap()
}

func meta_drop(m *Meta) {
	m.Drop()
	if m.Depth() < 1 {
		m.Put(NewSystemStack())
	}
}

func meta_new_system_stack(m *Meta) {
	s := NewSystemStack()
	s.Push(m.Current())
	m.Put(s)
}

func meta_swap(m *Meta) {
	m.Swap()
}

// CURRENT STACK

func current_decap(m *Meta) {
	m.Current().Decap()
}

func current_depth(m *Meta) {
	m.Current().Push(N(m.Current().Depth()))
}

func current_drop(m *Meta) {
	m.Current().Drop()
}

func current_dup(m *Meta) {
	m.Current().Dup()
}

func current_over(m *Meta) {
	m.Current().Over()
}

func current_purge(m *Meta) {
	m.Current().Purge()
}

func current_rot(m *Meta) {
	m.Current().Rot()
}

func current_swap(m *Meta) {
	m.Current().Swap()
}

// CONCURRENCY

func bkg(m *Meta) {
	block := m.Current().Pop().(B)

	items := NewStack("bkg")
	items.Push(m.Current().Pop())

	coBC("bkg", items, block)
}

func co(m *Meta) {
	filename := string(m.Current().Pop().(T))

	in := NewQueue()
	out := NewQueue()

	stack := NewStack("co")
	stack.Push(in)
	stack.Push(out)

	code := loadFile(filename)
	tokens := parse(code)
	file_bc := compile(tokens)

	m.Current().Push(out)
	m.Current().Push(in)

	coBC("co", stack, file_bc)
}

func work(m *Meta) {
	block := m.Current().Pop().(B)
	in := m.Current().Items[len(m.Current().Items)-1].(*Queue)
	out := m.Current().Items[len(m.Current().Items)-2].(*Queue)

	stack := NewStack("work")
	stack.Push(in)
	stack.Push(out)

	m.Current().Push(out)
	m.Current().Push(in)

	coBC("work", stack, block)
}

func wait(m *Meta) {
	threads.Wait()
}

// DEBUG

func bl_println(m *Meta) {
	print(m.Current().Pop().Print(), "\n")
}

func bl_refl(m *Meta) {
	m.Current().Push(T(m.Current().Pop().Refl()))
}

func bl_warn(m *Meta) {
	NOPE("warn")
}

// LOADING

func do(m *Meta) {
	filename := string(m.Current().Pop().(T))

	code := loadFile(filename)
	tokens := parse(code)
	file_bc := compile(tokens)

	doBC(m, file_bc)
}

func bload(m *Meta) {
	filename := string(m.Current().Pop().(T))

	code := loadFile(filename)
	tokens := parse(code)
	file_bc := compile(tokens)
	m.Current().Push(B(file_bc))
}

func imp(m *Meta) {
	NOPE("imp")
}

// MATH

func add(m *Meta) {
	n1 := m.Current().Pop().(N)
	n2 := m.Current().Pop().(N)

	m.Current().Push(n2 + n1)
}

func sub(m *Meta) {
	n1 := m.Current().Pop().(N)
	n2 := m.Current().Pop().(N)

	m.Current().Push(n2 - n1)
}

func div(m *Meta) {
	n1 := m.Current().Pop().(N)
	n2 := m.Current().Pop().(N)

	m.Current().Push(n2 / n1)
}

func mod(m *Meta) {
	n1 := m.Current().Pop().(N)
	n2 := m.Current().Pop().(N)

	m.Current().Push(n2 % n1)
}

func mul(m *Meta) {
	n1 := m.Current().Pop().(N)
	n2 := m.Current().Pop().(N)

	m.Current().Push(n2 * n1)
}

func n_to_r(m *Meta) {
	m.Current().Push(m.Current().Pop().(N).N_to_R())
}

func n_to_t(m *Meta) {
	m.Current().Push(m.Current().Pop().(N).N_to_T())
}

// IO

func read(m *Meta) {
	source := m.Current().Pop()
	q := m.Current().Peek().(*Queue)
	io := ReadIO(source, q)
	m.Current().Push(io)
}

func write(m *Meta) {
	dest := m.Current().Pop()
	q := m.Current().Peek().(*Queue)
	io := WriteIO(dest, q)
	m.Current().Push(io)
}

// LOGIC AND LOOPS

func bl_either(m *Meta) {
	comp := m.Current().Pop().(B)
	iffalse := m.Current().Pop().(B)
	iftrue := m.Current().Pop().(B)
	doBC(m, comp)
	if m.Current().Pop().(*Tag).Kind == "true" {
		doBC(m, iftrue)
	} else {
		doBC(m, iffalse)
	}
}

func bl_eq(m *Meta) {
	i1 := m.Current().Pop()
	i2 := m.Current().Pop()
	m.Current().Push(blEq(i1, i2))
}

func bl_if(m *Meta) {
	comp := m.Current().Pop().(B)
	actn := m.Current().Pop().(B)
	doBC(m, comp)
	if m.Current().Pop().(*Tag).Kind == "true" {
		doBC(m, actn)
	}
}

func bl_is(m *Meta) {
	NOPE("is")
}

func bl_not(m *Meta) {
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
}

func bl_until(m *Meta) {
	comp := m.Current().Pop().(B)
	actn := m.Current().Pop().(B)
Until:
	for {
		doBC(m, comp)
		if m.Current().Pop().(*Tag).Kind == "true" {
			break Until
		}
		doBC(m, actn)
	}
}

func bl_while(m *Meta) {
	comp := m.Current().Pop().(B)
	actn := m.Current().Pop().(B)
While:
	for {
		doBC(m, comp)
		if m.Current().Pop().(*Tag).Kind != "true" {
			break While
		}
		doBC(m, actn)
	}
}

func bl_loop(m *Meta) {
	actn := m.Current().Pop().(B)
	for {
		doBC(m, actn)
	}
}

// OBJECTS

func o_new(m *Meta) {
	m.Current().Push(NewObject())
}

func o_self(m *Meta) {
	m.Self()
}

func o_child(m *Meta) {
	o := m.Object()
	child := NewChildObject(o)
	m.Current().Push(child)
}

func o_fetch(m *Meta) {
	slot := m.Current().Pop().(W)
	o := m.Object()
	m.Current().Push(o.Fetch(slot))
}

func o_get(m *Meta) {
	slot := m.Current().Pop().(W)
	m.Object().Get(m, slot)
}

func o_set(m *Meta) {
	slot := m.Current().Pop().(W)
	i := m.Current().Pop()
	m.Object().Set(slot, i)
}

// QUEUES

func q_new(m *Meta) {
	m.Current().Push(NewQueue())
}

func q_deq(m *Meta) {
	i := m.Current().Peek().(*Queue).Dequeue()
	m.Current().Push(i)
}

func q_enq(m *Meta) {
	i := m.Current().Pop()
	m.Current().Peek().(*Queue).Enqueue(i)
}

func q_proq(m *Meta) {
	block := m.Current().Pop().(B)
	q := m.Current().Pop().(*Queue)

ProcQLoop:
	for {
		select {
		case item := <-q.Items:
			m.Current().Push(item)
			doBC(m, block)
		default:
			break ProcQLoop
		}
	}
}

func q_to_s(m *Meta) {
	NOPE("q-to-s")
}

func q_to_v(m *Meta) {
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

	m.Current().Push(V(items))
}

func q_to_t(m *Meta) {
	q := m.Current().Pop().(*Queue)
	text := T("")

QtoT:
	for {
		i := <-q.Items
		if blEq(i, NewNil("q_to_t")).Bool {
			break QtoT
		} else {
			text = text.App(i).(T)
		}
	}

	m.Current().Push(text)
}

func q_unq(m *Meta) {
	NOPE("unq")
}

// STACKS

func s_new(m *Meta) {
	m.Current().Push(NewStack("user"))
}

func s_pop(m *Meta) {
	m.Current().Push(m.Current().Peek().(stackable).Pop())
}

func s_push(m *Meta) {
	d := m.Current().Pop()
	m.Current().Peek().(stackable).Push(d)
}

func s_size(m *Meta) {
	m.Current().Push(N(m.Current().Peek().(stackable).Depth()))
}

func s_tail(m *Meta) {
	m.Current().Peek().(stackable).Decap()
}

// VECTORS AND SEQUENCES

func v_new(m *Meta) {
	m.Current().Push(V{})
}

func seq_cat(m *Meta) {
	i1 := m.Current().Pop().(sequence)
	i2 := m.Current().Pop().(sequence)
	result := i2.Cat(i1)
	m.Current().Push(result)
}

func seq_app(m *Meta) {
	i := m.Current().Pop()
	v := m.Current().Pop().(sequence)
	m.Current().Push(v.App(i))
}

func seq_ato(m *Meta) {
	c := m.Current()
	n := c.Pop().(N)
	v := m.Current().Peek().(sequence)
	i := v.Ato(n)
	m.Current().Push(i)
}

func seq_rmo(m *Meta) {
	n := m.Current().Pop().(N)
	v := m.Current().Pop().(sequence)
	nv := v.Rmo(n)
	m.Current().Push(nv)
}

func seq_len(m *Meta) {
	v := m.Current().Peek().(sequence)
	m.Current().Push(N(v.Len()))
}

func seq_pick(m *Meta) {
	NOPE("pick")
}
func seq_del(m *Meta) {
	NOPE("del")
}
func seq_emt(m *Meta) {
	NOPE("emt")
}
func v_to_s(m *Meta) {
	NOPE("v-to-s")
}
func v_to_q(m *Meta) {
	NOPE("v-to-q")
}

// BLOCK

func block_call(m *Meta) {
	doBC(m, m.Current().Pop().(B))
}

func block_decompile(m *Meta) {
	b := m.Current().Pop().(B)
	t := analyze(b)
	m.Current().Push(t)
}

func block_disassemble(m *Meta) {
	b := m.Current().Pop().(B)
	v := b.Disassemble()
	m.Current().Push(v)
}

// TEXT

func bl_print(m *Meta) {
	print(m.Current().Pop().(T))
}

func t_to_cv(m *Meta) {
	cv := m.Current().Pop().(T).T_to_CV()
	m.Current().Push(cv)
}

func t_to_b(m *Meta) {
	code := []rune(string(m.Current().Pop().(T)))
	tokens := parse(code)
	ops := compile(tokens)
	b := B(ops)
	m.Current().Push(b)
}

// tags
func tag_to_t(m *Meta) {
	NOPE("?-to-t")
}

func bl_true(m *Meta) {
	m.Current().Push(NewTrue("true"))
}

func bl_nil(m *Meta) {
	m.Current().Push(NewNil("nil"))
}

// RUNES

func r_to_t(m *Meta) {
	m.Current().Push(m.Current().Pop().(R).R_to_T())
}

func r_to_n(m *Meta) {
	m.Current().Push(m.Current().Pop().(R).R_to_N())
}

// OCTETS

func b_and(m *Meta) {
	c1 := m.Current().Pop().(C)
	c2 := m.Current().Pop().(C)
	c3 := c1.Band(c2)
	m.Current().Push(c3)
}

func b_or(m *Meta) {
	c1 := m.Current().Pop().(C)
	c2 := m.Current().Pop().(C)
	c3 := c1.Bor(c2)
	m.Current().Push(c3)
}

func b_xor(m *Meta) {
	c1 := m.Current().Pop().(C)
	c2 := m.Current().Pop().(C)
	c3 := c1.Bxor(c2)
	m.Current().Push(c3)
}

func b_shift_l(m *Meta) {
	c1 := m.Current().Pop().(C)
	c2 := m.Current().Pop().(C)
	c3 := c1.Band(c2)
	m.Current().Push(c3)
}

func b_shift_r(m *Meta) {
	c1 := m.Current().Pop().(C)
	c2 := m.Current().Pop().(C)
	c3 := c1.Band(c2)
	m.Current().Push(c3)
}

func c_to_r(m *Meta) {
	m.Current().Push(m.Current().Pop().(C).C_to_R())
}

func c_to_n(m *Meta) {
	m.Current().Push(m.Current().Pop().(C).C_to_N())
}
