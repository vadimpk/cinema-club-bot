package event

type Bus interface {
	// Publish is used to publish a new event with topic and data.
	Publish(topic string, event interface{})
	// Subscribe is used to subscribe to event topic and add handler.
	Subscribe(topic string, handler interface{}) error
}
