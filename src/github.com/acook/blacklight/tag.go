package main

type Tag struct {
	Kind  string
	Label string
	Data  interface{}
	Bool  bool
}

func (t Tag) String() string {
	return t.Kind + "#" + t.Label
}

func (t Tag) Value() interface{} {
	return t
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

func NewFDTag(label string, handle *FD) *Tag {
	t := NewTag("FD", label)
	t.Data = handle
	return t
}

func NewFileTag(label string, handle *FD) *Tag {
	t := NewTag("FILE", label)
	t.Data = handle
	return t
}

func NewErr(label string, data datatypes) *Tag {
	t := NewTag("ERR", label)
	t.Data = data
	return t
}
