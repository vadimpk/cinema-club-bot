package event

import (
	messagebus "github.com/vardius/message-bus"
)

var _ Bus = (*msgBus)(nil)

type msgBus struct {
	bus messagebus.MessageBus
}

func NewMsgBus(bufferSize int) *msgBus {
	return &msgBus{
		bus: messagebus.New(bufferSize),
	}
}

func (b *msgBus) Publish(topic string, event interface{}) {
	b.bus.Publish(topic, event)
}

func (b *msgBus) Subscribe(topic string, handler interface{}) error {
	return b.bus.Subscribe(topic, handler)
}
