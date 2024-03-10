package main

import (
	"fmt"
	"reflect"
)

type Object struct {
	Slots  map[W]datatypes
	Parent *Object
}

func NewObject() *Object {
	o := new(Object)
	o.Slots = make(map[W]datatypes)
	return o
}

func NewChildObject(parent *Object) *Object {
	o := NewObject()
	o.Parent = parent
	return o
}

func (o *Object) Set(w W, i datatypes) {
	o.Slots[w] = i
}

func (o *Object) Fetch(w W) datatypes {
	i, found := o.Slots[w]
	if !found {
		panic("Object.Fetch: slot `" + w.Refl() + "` does not exist!")
	}
	return i
}

func (o *Object) Get(meta *Meta, w W) {
	meta.ObjectStack.Push(o)
	defer meta.ObjectStack.Pop()

	ok := o.DeleGet(meta, w)

	if !ok {
		print("\n")
		print("error in Object.Get: ")
		print("slot `", w.Refl(), "` not found!\n")
		print(" -- given: ", w.Refl(), "\n")
		print(" --   has: ", o.Labels().Refl(), "\n")
		panic("Object.Get: slot `" + w.Refl() + "` does not exist!")
	}
}

func (o *Object) DeleGet(meta *Meta, w W) bool {
	current := meta.Current()
	i, found := o.Slots[w]

	if found {
		switch i.(type) {
		case B:
			doBC(meta, i.(B))
		default:
			current.Push(i)
		}
	} else if o.Parent != nil {
		found = o.Parent.DeleGet(meta, w)
	}

	return found
}

func (o Object) Refl() string {
	str := ""
	
	list := V{}
	_, s := o.DeepRefl(0, list)
	str += s

	return str + ""
}

func (o Object) ShallowRefl() string {
	str := "\033[32mO{ "

	for k, v := range o.Slots {
		str += k.Refl() + ":" 

		switch t := v.(type) {
		case *Object:
			str += v.(*Object).SimpleRefl()
		case N, C, T, B, OP:
			str += v.Refl()
		default:
			str += "..."
			str += fmt.Sprint(reflect.TypeOf(t))
			str += "..."
		}

		str += " "
	}

	return str + "}\033[0m"
}

func (o Object) DeepRefl(depth N, list V) (V, string) {
	oc := func(n N) string {
		c := "\033[38;5;"
		c += cycle([]string{"46", "34", "28", "22"}, n)
		c += "m"
		return c
	}

	str := oc(depth)
	str += "O{ "

	for k, v := range o.Slots {
		str += k.Refl() + ":" 

		switch t := v.(type) {
		case N, C, T, B, OP:
			str += oc(depth + 1)
			str += v.Refl()
			str += oc(depth)
		default:
			if list.Contains(v).Kind == "true" {
				str += oc(depth + 1)
				str += "..."
				str += fmt.Sprint(reflect.TypeOf(t))
				str += "..."
				str += oc(depth)
			} else {
				list = list.App(v).(V)
				l, s := v.DeepRefl(depth + 1, list)
				list = l
				str += s
			}
		}

		str += " " + oc(depth)
	}

	return list, str + "}\033[0m"
}

func cycle(v []string, n N) string {
	return v[int(n.Value().(int64)) % len(v)]
}

func (o Object) SimpleRefl() string {
	str := "O#"

	str = str + fmt.Sprint(len(o.Slots))

	return str + ""
}

func (o Object) Value() interface{} {
	return o.Slots
}

func (o Object) Labels() V {
	var labels V
	for label := range o.Slots {
		labels = append(labels, label)
	}
	return labels
}
