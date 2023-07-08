package pubsub

import (
	"testing"
)

func TestPubSub(t *testing.T) {
	ps := New()

	// Test subscribing and publishing to a topic
	topic := "testTopic"
	handlerCalled := false
	handler := func(args ...any) {
		handlerCalled = true
	}

	err := ps.Subscribe(topic, handler)
	if err != nil {
		t.Errorf("Subscribe returned an error: %s", err.Error())
	}

	err = ps.Publish(topic, "test message")
	if err != nil {
		t.Errorf("Publish returned an error: %s", err.Error())
	}

	if !handlerCalled {
		t.Error("Handler was not called after publishing")
	}

	// Test unsubscribing
	handlerCalled = false

	err = ps.Unsubscribe(topic)
	if err != nil {
		t.Errorf("Unsubscribe returned an error: %s", err.Error())
	}

	err = ps.Publish(topic, "test message")
	if err != nil {
		t.Errorf("Publish returned an error: %s", err.Error())
	}

	if handlerCalled {
		t.Error("Handler was called after unsubscribing")
	}

	// Test closing a topic
	topic = "testTopic2"
	handlerCalled = false

	err = ps.Subscribe(topic, handler)
	if err != nil {
		t.Errorf("Subscribe returned an error: %s", err.Error())
	}

	err = ps.CloseTopic(topic)
	if err != nil {
		t.Errorf("CloseTopic returned an error: %s", err.Error())
	}

	err = ps.Publish(topic, "test message")
	if err != nil {
		t.Errorf("Publish returned an error: %s", err.Error())
	}

	if handlerCalled {
		t.Error("Handler was called after closing the topic")
	}

	// Test shutting down the PubSub
	topic = "testTopic3"
	handlerCalled = false

	err = ps.Subscribe(topic, handler)
	if err != nil {
		t.Errorf("Subscribe returned an error: %s", err.Error())
	}

	err = ps.Shutdown()
	if err != nil {
		t.Errorf("Shutdown returned an error: %s", err.Error())
	}

	err = ps.Publish(topic, "test message")
	if err != nil {
		t.Errorf("Publish returned an error: %s", err.Error())
	}

	if handlerCalled {
		t.Error("Handler was called after shutting down the PubSub")
	}

	// Test unsubscribing from a non-existent topic
	err = ps.Unsubscribe("nonExistentTopic")
	if err != nil {
		t.Errorf("Unsubscribe returned an error for a non-existent topic: %s", err.Error())
	}

	// Test closing a non-existent topic
	err = ps.CloseTopic("nonExistentTopic")
	if err != nil {
		t.Errorf("CloseTopic returned an error for a non-existent topic: %s", err.Error())
	}

	// Test publishing to a non-existent topic
	err = ps.Publish("nonExistentTopic", "test message")
	if err != nil {
		t.Errorf("Publish returned an error for a non-existent topic: %s", err.Error())
	}
}
