package main

import (
	"strconv"
	"sync"
)

type Meta struct {
	sync.Mutex
	Id          int
	Items       []*Stack
	ObjectStack *ObjectStack
	SelfFlag    bool
}

func NewMeta() *Meta {
	m := &Meta{}
	m.Id = getStackId()
	m.ObjectStack = NewObjectStack()
	return m
}

func (m Meta) Value() interface{} {
	return m
}

func (m Meta) String() string {
	str := "$" + strconv.Itoa(m.Id) + "#" + strconv.Itoa(m.Depth()) + "< "

	for _, i := range m.Items {
		str += i.String() + " "
	}

	return str + ">"
}

// for stack interface compatibility

func (m *Meta) Push(i datatypes) {
	m.Lock()
	defer m.Unlock()

	m.Items = append(m.Items, i.(*Stack))
}

func (m *Meta) Pop() datatypes { // quite dangerous
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

	s := m.Items[m.Depth()-1]
	m.Items = m.Items[:m.Depth()-1]
	return s
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
	m.Lock()
	defer m.Unlock()

	if m.Depth() > 0 {
		m.Items = m.Items[:m.Depth()-1]
	}
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
