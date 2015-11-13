package main

type W uint64

func (w W) String() string {
	return string(wd_map[uint64(w)])
}

func (w W) Value() interface{} {
	return uint64(w)
}

func (w W) Text() T {
	return T(wd_map[uint64(w)])
}
