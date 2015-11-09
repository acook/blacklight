package main

type ops_fifo struct {
	items [][]operation
}

func (f *ops_fifo) push(ops []operation) {
	f.items = append(f.items, ops)
}

func (f *ops_fifo) pop() []operation {
	ops := f.items[f.depth()-1]
	f.items = f.items[:f.depth()-1]
	return ops
}

func (f *ops_fifo) depth() int {
	return len(f.items)
}
