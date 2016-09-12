package main

// functions that are called on datatypes interface
// rather than being methods on a datatype

func blEq(i1 datatypes, i2 datatypes) *Tag {
	if i1.Value() == i2.Value() {
		return NewTrue("eq")
	} else {
		return NewNil("eq")
	}
}
