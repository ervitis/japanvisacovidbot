package queue

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func New() *gobbit {
	return &gobbit{
		subs: make(map[string]chan *Message),
	}
}

func (g *gobbit) CreateTopic(topic string) error {
	_, ok := g.subs[topic]
	if ok {
		return fmt.Errorf("topic %s already exists", topic)
	}

	g.subs[topic] = make(chan *Message, 5)
	return nil
}

func (g *gobbit) Publish(topic string, data interface{}) {
	_, ok := g.subs[topic]
	if !ok {
		panic(fmt.Errorf("topic %s does not exist", topic))
	}

	metadata := Metadata{
		Timestamp: time.Now(),
	}

	g.l.Lock()
	metadata.ID = uuid.New().String()
	g.order += 1
	g.l.Unlock()

	message := &Message{
		Payload:  data,
		MetaData: metadata,
	}

	go func() {
		g.subs[topic] <- message
	}()

}

func (g *gobbit) Subscribe(topic string, handler func(message *Message)) {
	_, ok := g.subs[topic]
	if !ok {
		panic(fmt.Errorf("topic %s does not exist", topic))
	}

	_, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
	}()

	go func() {
		for {
			for _, sub := range g.subs {
				select {
				case msg := <-sub:
					handler(msg)
				}
			}
		}
	}()
}
