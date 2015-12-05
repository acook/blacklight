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
// 0xF3 : rune - uint32 (rune)
// 0xF4 : number - int64
// 0xF5 : float - float64
// 0xF6 : text - length:uint64 data:Rs
// 0xF7 : block - length:uint64 data:bc
// 0xF8 : vector - data:bc
// 0xF9 : endvector
// 0xFA : tag - kind:uint8 metadata:uint32 msg:text
// 0xFB-0xFE : FUTURE DATATYPES
// 0xFF : RESERVED EXTENDED FLAG

var inb_map = map[string]byte{ // item name->byte map
	"opword":    0xF0,
	"word":      0xF1,
	"byte":      0xF2,
	"rune":      0xF3,
	"number":    0xF4,
	"float":     0xF5,
	"text":      0xF6,
	"block":     0xF7,
	"vector":    0xF8,
	"endvector": 0xF9,
	"tag":       0xFA,
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
		"say":   bl_println,
		"print": bl_print,
		"refl":  bl_refl,
		"warn":  bl_warn,

		// loading
		"do":    do,
		"imp":   imp,
		"bload": bload,

		// math
		"add":    add,
		"sub":    sub,
		"div":    div,
		"mod":    mod,
		"mul":    mul,
		"n-to-r": n_to_r,
		"n-to-t": n_to_t,

		// file io
		"read":  read,
		"write": write,

		// logic & loops
		"either": bl_either,
		"eq":     bl_eq,
		"if":     bl_if,
		"is":     bl_is,
		"not":    bl_not,
		"until":  bl_until,
		"while":  bl_while,
		"loop":   bl_loop,

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
		"del":    seq_del,
		"emt":    seq_emt,
		"pick":   seq_pick,
		"v-to-s": v_to_s,
		"v-to-q": v_to_q,

		// blocks
		"call":      block_call,
		"decompile": block_decompile,
		"dis":       block_disassemble,

		// text
		"t-to-cv": t_to_cv,
		"compile": t_to_b,

		// tags
		"?-to-t": tag_to_t,
		"true":   bl_true,
		"nil":    bl_nil,

		// runes
		"r-to-t": r_to_t,
		"r-to-n": r_to_n,

		// octets
		"band":    b_and,
		"bor":     b_or,
		"bxor":    b_xor,
		"bshiftl": b_shift_l,
		"bshiftr": b_shift_r,
		"c-to-r":  c_to_r,
		"c-to-n":  c_to_n,
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
