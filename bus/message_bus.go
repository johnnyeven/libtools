package bus

import (
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
	"github.com/sirupsen/logrus"
)

type MessageBus struct {}

func (b *MessageBus) Init() {
	// configure id generator (it doesn't have to be monoton)
	node := uint(1)
	initialTime := uint(0)
	if err := monoton.Configure(sequencer.NewMillisecond(), node, initialTime); err != nil {
		logrus.Panicf("[MessageBus] monoton.Configure err: %v", err)
	}

	// configure bus
	if err := bus.Configure(bus.Config{Next: monoton.Next}); err != nil {
		logrus.Panicf("[MessageBus] bus.Configure err: %v", err)
	}
}

func (b *MessageBus) RegisterTopic(topicNames ...string) {
	bus.RegisterTopics(topicNames...)
}

func (b *MessageBus) RegisterHandler(key, matcher string, handlerFunc func(e *bus.Event)) () {
	handler := bus.Handler{
		Handle:  handlerFunc,
		Matcher: matcher,
	}

	bus.RegisterHandler(key, &handler)
}

func (b *MessageBus) DeregisterHandler(key string) {
	bus.DeregisterHandler(key)
}

func (b *MessageBus) ListHandlerKeys() []string {
	return bus.ListHandlerKeys()
}

func (b *MessageBus) Emit(topicName string, data interface{}, txID string) (event *bus.Event, err error) {
	return bus.Emit(topicName, data, txID)
}
