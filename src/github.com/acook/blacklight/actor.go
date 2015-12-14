package main

import (
	"sync"
	"sync/atomic"
)

var ActorLock sync.Mutex
var Actors uint64

type actor_msgs_box []interface{}
type actor_ivar_map map[string]interface{}
type actor_messages chan actor_msgs_box
type actor_response chan actor_msgs_box

type Actor struct {
	id       uint64
	messages actor_messages
	ivars    actor_ivar_map
}

func ActorId() uint64 {
	return atomic.AddUint64(&Actors, 1)
}

func NewActor(buf uint8) *Actor {
	a := new(Actor)
	a.id = ActorId()
	print(a.id)
	a.messages = make(chan actor_msgs_box, buf)
	return a
}

func (a *Actor) Send(label string, args actor_msgs_box) actor_response {
	responder := make(actor_response, 1)
	message := actor_msgs_box{label, responder}
	a.messages <- append(message, args...)
	return responder
}

func (a *Actor) Trigger(label string, args actor_msgs_box) {
	message := actor_msgs_box{label, nil}
	a.messages <- append(message, args...)
}
