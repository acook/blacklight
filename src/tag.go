package main

type Tag struct {
	Kind  string
	Label string
	Data  interface{}
	Bool  bool
}

func (t Tag) Error() string {
	out := t.Print()

	switch t.Data.(type) {
	case sequence:
		out = out + "\n" + t.Data.(sequence).Print()
	case *IO:
		out = out + "<IO>"
	}

	return out
}

func (t Tag) Refl() string {
	return t.Print()
}

func (t Tag) Print() string {
	return t.Kind + "#" + t.Label
}

func (t Tag) Value() interface{} {
	return t.Kind
}

func (t Tag) Bytes() []byte {
	if t.Kind != "nil" {
		panic("Tag.Bytes: Attempt to serialize non-nil Tag!")
	}
	return nil
}

func (t Tag) Bytecode() []byte {
	return t.Bytes()
}

func NewTag(kind string, label string) *Tag {
	t := new(Tag)
	t.Kind = kind
	t.Label = label
	return t
}

func NewTrue(label string) *Tag {
	t := NewTag("true", label)
	t.Bool = true
	return t
}

func NewNil(label string) *Tag {
	t := NewTag("nil", label)
	return t
}

func NewFDTag(label string, handle *IO) *Tag {
	t := NewTag("FD", label)
	t.Data = handle
	return t
}

func NewFileTag(label string, handle *IO) *Tag {
	t := NewTag("FILE", label)
	t.Data = handle
	return t
}

func NewErr(label string, data datatypes) *Tag {
	t := NewTag("ERR", label)
	t.Data = data
	return t
}

func blOr(t1 *Tag, t2 *Tag) *Tag {
	if t1.Bool || t2.Bool {
		return NewTrue("or")
	}
	return NewNil("or")
}

func blAnd(t1 *Tag, t2 *Tag) *Tag {
	if t1.Bool && t2.Bool {
		return NewTrue("and")
	}
	return NewNil("and")
}
