package main

import (
	"sync"
	"sync/atomic"
)

var actorLock sync.Mutex
var actors uint64

type actor_msgs_box []datatypes
type actor_ivar_map map[W]datatypes
type actor_messages chan actor_msgs_box

type Actor struct {
	id       uint64
	messages actor_messages
	ivars    actor_ivar_map
	meta     *Meta
}

func actorID() uint64 {
	return atomic.AddUint64(&actors, 1)
}

func NewActor() *Actor {
	a := new(Actor)
	a.id = actorID()
	print(a.id)
	a.messages = make(chan actor_msgs_box, 8)
	meta := NewMeta()
	meta.Push(NewStack("ACT#" + N(a.id).Print()))
	go a.Perform() // start the actor receiving messages
	return a
}

func (a *Actor) Send(label W, args V) *Queue {
	responder := new(Queue)
	message := actor_msgs_box{label, responder, args}
	a.messages <- message
	return responder
}

func (a *Actor) Trigger(label W, args V) {
	message := actor_msgs_box{label, nil, args}
	a.messages <- message
}

func (a *Actor) Perform() {
	for {
		msg := <-a.messages

		label := msg[0].(W)
		resp := msg[1].(Queue)
		args := msg[2].(V)

		i, found := a.ivars[label]

		if found {
			switch i.(type) {
			case B:
				a.meta.Put(NewStack("Actor.Perform"))
				a.meta.Current().Items = append(a.meta.Current().Items, args...)
				doBC(a.meta, i.(B))
				resp.Enqueue(a.meta.Current().S_to_V())
			default:
				resp.Enqueue(i)
			}
		} else {
			panic("Actor slot not found")
		}

	}
}

// for datatypes interface compliance

func (a Actor) Print() string {
	return a.Refl()
}

func (a Actor) Refl() string {
	str := "ACT#" + N(a.id).Print() + "<"

	for k, v := range a.ivars {
		str += k.Print() + ":" + v.Print() + " "
	}

	return str + ">"
}

func (a Actor) Value() interface{} {
	return a
}
