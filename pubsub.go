package pubsub

type Operation int

const (
	Subscribe Operation = iota
	SubscribeOnce
	SubscribeOnceEach
	Publish
	TryPublish
	Unsubscribe
	UnsubscribeAll
	CloseTopic
	Shutdown
)

// Message is a message sent to a pubsub instance.
// It contains the operation to perform and the arguments to pass to the operation.
type Message struct {
	Topic     string
	Operation Operation
	Args      []any
}

// Subscriber is the interface that wraps the Subscribe, SubscribeOnce, SubscribeOnceEach, Unsubscribe and UnsubscribeAll methods.
// Subscribe adds a handler to the topic.
// SubscribeOnce adds a handler to the topic and removes it after the first call.
// SubscribeOnceEach adds a handler to the topic and removes it after the first call for each handler.
// Unsubscribe removes all handlers from the topic.
// UnsubscribeAll removes all handlers from all topics
type Subscriber interface {
	Subscribe(topic string, handler func(...any)) error
	SubscribeOnce(topic string, handler func(...any)) error
	SubscribeOnceEach(topic string, handler func(...any)) error
	Unsubscribe(topic string) error
	UnsubscribeAll() error
}

// Publisher is the interface that wraps the Publish and TryPublish methods.
// Publish calls all handlers for the topic.
// TryPublish calls all handlers for the topic and returns the first error.
type Publisher interface {
	Publish(topic string, args ...any) error
	TryPublish(topic string, args ...any) error
}

// PubSub is the interface that groups the Subscriber and Publisher interfaces.
// It also adds the CloseTopic and Shutdown methods.
// CloseTopic removes all handlers from the topic and deletes the topic.
// Shutdown removes all handlers from all topics and deletes all topics.
type PubSub interface {
	Subscriber
	Publisher
	CloseTopic(topic string) error
	Shutdown() error
}

// New returns a new PubSub instance.
func New() PubSub {
	return &pubsub{
		topics: make(map[string]*topic),
	}
}

type pubsub struct {
	topics map[string]*topic
}

// Shutdown removes all handlers from all topics and deletes all topics.
func (p *pubsub) Subscribe(topic string, handler func(...any)) error {
	return p.subscribe(topic, handler, false, false)
}

// SubscribeOnce adds a handler to the topic and removes it after the first call.
func (p *pubsub) SubscribeOnce(topic string, handler func(...any)) error {
	return p.subscribe(topic, handler, true, false)
}

// SubscribeOnceEach adds a handler to the topic and removes it after the first call for each handler.
func (p *pubsub) SubscribeOnceEach(topic string, handler func(...any)) error {
	return p.subscribe(topic, handler, true, true)
}

// CloseTopic removes all handlers from the topic and deletes the topic.
func (p *pubsub) subscribe(topic string, handler func(...any), once, onceEach bool) error {
	t, ok := p.topics[topic]
	if !ok {
		t = newTopic(topic)
		p.topics[topic] = t
	}
	return t.subscribe(handler, once, onceEach)
}

// CloseTopic removes all handlers from the topic and deletes the topic.
func (p *pubsub) Unsubscribe(topic string) error {
	t, ok := p.topics[topic]
	if !ok {
		return nil
	}
	return t.unsubscribe()
}

// UnsubscribeAll removes all handlers from all topics.
func (p *pubsub) UnsubscribeAll() error {
	for _, t := range p.topics {
		if err := t.unsubscribe(); err != nil {
			return err
		}
	}
	return nil
}

// Publish calls all handlers for the topic.
func (p *pubsub) Publish(topic string, args ...any) error {
	return p.publish(topic, args, false)
}

// TryPublish calls all handlers for the topic and returns the first error.
func (p *pubsub) TryPublish(topic string, args ...any) error {
	return p.publish(topic, args, true)
}

func (p *pubsub) publish(topic string, args []any, try bool) error {
	t, ok := p.topics[topic]
	if !ok {
		return nil
	}
	return t.publish(args, try)
}

// CloseTopic removes all handlers from the topic and deletes the topic.
func (p *pubsub) CloseTopic(topic string) error {
	t, ok := p.topics[topic]
	if !ok {
		return nil
	}
	return t.close()
}

// Shutdown removes all handlers from all topics and deletes all topics.
func (p *pubsub) Shutdown() error {
	for _, t := range p.topics {
		if err := t.close(); err != nil {
			return err
		}
	}
	return nil
}

type topic struct {
	name     string
	handlers []func(...any)
	once     bool
	onceEach bool
	closed   bool
}

func newTopic(name string) *topic {
	return &topic{
		name: name,
	}
}

func (t *topic) subscribe(handler func(...any), once, onceEach bool) error {
	if t.closed {
		return nil
	}
	t.handlers = append(t.handlers, handler)
	t.once = once
	t.onceEach = onceEach
	return nil
}

func (t *topic) unsubscribe() error {
	if t.closed {
		return nil
	}
	t.handlers = nil
	t.closed = true
	return nil
}

func (t *topic) publish(args []any, try bool) error {
	if t.closed {
		return nil
	}
	for _, handler := range t.handlers {
		handler(args...)
		if t.once {
			t.handlers = nil
		}
		if t.onceEach {
			handler = nil
		}
	}
	return nil
}

func (t *topic) close() error {
	if t.closed {
		return nil
	}
	t.handlers = nil
	t.closed = true
	return nil
}
