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
// 0xF6 : text - length:uint64 data:Cs
// 0xF7 : block - length:uint64 data:bc
// 0xF8 : vector - length:uint64 data:items
// 0xF9 : tag - kind:uint8 metadata:uint32 msg:text
// 0xFB-0xFE : FUTURE DATATYPES
// 0xFF : RESERVED EXTENDED FLAG

var inb_map = map[string]byte{ // item name->byte map
	"opword": 0xF0,
	"word":   0xF1,
	"byte":   0xF2,
	"char":   0xF3,
	"number": 0xF4,
	"float":  0xF5,
	"text":   0xF6,
	"block":  0xF7,
	"vector": 0xF8,
	"tag":    0xF9,
}

var op_map map[string]byte
var fn_map map[byte]func(m *Meta)
var lk_map map[byte]string
var total_ops uint8

func prepare_op_table() {

	var op_fn_map = map[string]func(*Meta){

		// meta
		"@":      push_current,
		"^":      push_last,
		"$":      push_meta,
		"$decap": meta_decap,
		"$drop":  meta_drop,
		"$new":   meta_new_system_stack,
		"$swap":  meta_swap,

		// current stack
		"decap": current_decap,
		"depth": current_depth,
		"drop":  current_drop,
		"dup":   current_dup,
		"over":  current_over,
		"purge": current_purge,
		"rot":   current_rot,
		"swap":  current_swap,

		// concurrency
		"bkg":  bkg,
		"co":   co,
		"work": work,
		"wait": wait,

		// debug
		"print": bl_print,
		"refl":  bl_refl,
		"warn":  bl_warn,

		// loading
		"do":  do,
		"imp": imp,

		// math
		"add":    add,
		"sub":    sub,
		"div":    div,
		"mod":    mod,
		"mul":    mod,
		"n-to-c": n_to_c,
		"n-to-t": n_to_t,

		// file io
		"read":  read,
		"write": write,

		// logic & loops
		"either": bl_either,
		"eq": func(m *Meta) {
			i1 := m.Current().Pop()
			i2 := m.Current().Peek()
			m.Current().Push(blEq(i1, i2))
		},
		"if":    bl_if,
		"is":    bl_is,
		"not":   bl_not,
		"until": bl_until,
		"while": bl_while,
		"loop":  bl_loop,

		// objects
		"o-new": o_new,
		"self":  o_self,
		"child": o_child,
		"fetch": o_fetch,
		"get":   o_get,
		"set":   o_set,

		// queues
		"q-new":  q_new,
		"deq":    q_deq,
		"enq":    q_enq,
		"proq":   q_proq,
		"q-to-s": q_to_s,
		"q-to-v": q_to_v,
		"q-to-t": q_to_t,
		"unq":    q_unq,

		// stacks
		"s-new": s_new,
		"pop":   s_pop,
		"push":  s_push,
		"size":  s_size,
		"tail":  s_tail,

		// Vectors & Sequences
		"v-new":  v_new,
		"cat":    seq_cat,
		"app":    seq_app,
		"ato":    seq_ato,
		"rmo":    seq_rmo,
		"len":    seq_len,
		"call":   block_call,
		"del":    seq_del,
		"emt":    seq_emt,
		"pick":   seq_pick,
		"v-to-s": v_to_s,
		"v-to-q": v_to_q,

		// tags
		"?-to-t": tag_to_t,
		"true":   bl_true,
		"nil":    bl_nil,

		// chars
		"c-to-t": c_to_t,
		"c-to-n": c_to_n,
	}

	op_map = make(map[string]byte)
	fn_map = make(map[byte]func(*Meta))
	lk_map = make(map[byte]string)

	var i byte = 1
	for k, v := range op_fn_map {
		op_map[k] = i
		fn_map[i] = v
		lk_map[i] = k
		i++
	}
	total_ops = uint8(i)

}
