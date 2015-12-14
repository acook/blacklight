package main

import (
	"sync"
	"sync/atomic"
)

var actorLock sync.Mutex
var actors uint64

type actor_msgs_box []datatypes
type actor_slot_map map[W]datatypes
type actor_messages chan actor_msgs_box

type Actor struct {
	id       uint64
	messages actor_messages
	slots    actor_slot_map
	meta     *Meta
}

func actorID() uint64 {
	return atomic.AddUint64(&actors, 1)
}

func NewActor() *Actor {
	a := new(Actor)

	a.id = actorID()
	a.slots = make(actor_slot_map)
	a.messages = make(chan actor_msgs_box, 8)
	meta := NewMeta()
	meta.Push(NewStack("ACT#" + N(a.id).Print()))

	// start the actor receiving messages concurrently
	go a.Perform()

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
	var has_block bool
	for {
		msg := <-a.messages

		label := msg[0].(W)
		resp := msg[1]
		args := msg[2].(V)

		i, found := a.slots[label]

		if found {
			switch i.(type) {
			case B:
				a.meta.Put(NewStack("Actor.Perform"))
				a.meta.Current().Items = append(a.meta.Current().Items, args...)
				doBC(a.meta, i.(B))
			default:
				has_block = true
			}

			switch resp.(type) {
			case Queue:
				q := resp.(Queue)
				if has_block {
					q.Enqueue(a.meta.Current().S_to_V())
				} else {
					q.Enqueue(i)
				}
				q.Close()
			}
		} else {
			print("\n")
			print("error in Actor.Send: ")
			print("slot `", label.Print(), "` not found!\n")
			print(" -- given: ", label.Print(), "\n")
			print(" --   has: ", a.Labels().Print(), "\n")
			panic("Object.Get: slot `" + label.Print() + "` does not exist!")
		}

	}
}

func (a *Actor) Labels() V {
	var labels V
	for label, _ := range a.slots {
		labels = append(labels, label)
	}
	return labels
}

// for datatypes interface compliance

func (a Actor) Print() string {
	return a.Refl()
}

func (a Actor) Refl() string {
	str := "ACT#" + N(a.id).Print() + "<"

	for k, v := range a.slots {
		str += k.Print() + ":" + v.Print() + " "
	}

	return str + ">"
}

func (a Actor) Value() interface{} {
	return a
}
