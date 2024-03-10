package main

import (
	"fmt"
	"strconv"
	"sync"
)

type Meta struct {
	sync.Mutex
	ID          uint64
	Items       []*Stack
	ObjectStack *ObjectStack
	SelfFlag    bool
}

func NewMeta() *Meta {
	m := &Meta{}
	m.ID = getStackID()
	m.ObjectStack = NewObjectStack()
	return m
}

func (m *Meta) Value() interface{} {
	return m
}

func (m *Meta) Refl() string {
	str := "$" + fmt.Sprint(m.ID) + "#" + strconv.Itoa(m.Depth()) + "< "

	for _, i := range m.Items {
		str += i.Refl() + " "
	}

	return str + ">"
}

func (m *Meta) DeepRefl(list V) (V, string) {
	return list, m.Refl()
}

// for stack interface compatibility
// will panic if you try to push anything other than a stack
func (m *Meta) Push(i datatypes) {
	m.Lock()
	defer m.Unlock()

	m.Items = append(m.Items, i.(*Stack))
}

// for stack interface compatibility
// quite dangerous because the stack could be in use
func (m *Meta) Pop() datatypes {
	m.Lock()
	defer m.Unlock()

	s := m.Items[m.Depth()-1]
	m.Items = m.Items[:m.Depth()-1]
	return s
}

// basic meta operations

func (m *Meta) Put(s *Stack) { // equivilent to push but typed
	m.Lock()
	defer m.Unlock()

	m.Items = append(m.Items, s)
}

func (m *Meta) Eject() *Stack { // equivilent to pop but typed, used internally
	m.Lock()
	defer m.Unlock()
	var s1 *Stack

	if m.Depth() > 0 {
		s1 = m.Items[m.Depth()-1]
		m.Items = m.Items[:m.Depth()-1]
	} else {
		s1 = NewSystemStack()
	}

	if m.Depth() < 1 {
		s2 := NewSystemStack()
		m.Items = append(m.Items, s2)
	}

	return s1
}

func (m *Meta) Peek() *Stack {
	m.Lock()
	defer m.Unlock()

	return m.Items[m.Depth()-1]
}

func (m *Meta) Depth() int {
	return len(m.Items)
}

func (m *Meta) Drop() {
	m.Eject()
}

func (m *Meta) Decap() {
	m.Lock()
	defer m.Unlock()

	depth := m.Depth()
	if depth > 1 {
		m.Items = m.Items[depth-1:]
	}
}

func (m *Meta) Swap() {
	m.Lock()
	defer m.Unlock()

	depth := m.Depth()
	if depth > 1 {
		m.Items = append(m.Items[:depth-2], m.Items[depth-1], m.Items[depth-2])
	}
}

// meta helpers

func (m *Meta) Current() *Stack {
	return m.Items[len(m.Items)-1]
}

func (m *Meta) Self() {
	m.Lock()
	defer m.Unlock()

	o := m.ObjectStack.Peek()
	m.Current().Push(o)
	m.SelfFlag = true
}

func (m *Meta) Object() *Object {
	return m.Current().Peek().(*Object)
}

func (m *Meta) NewStack(label string) {
	m.Put(NewStack(label))
}
