package message_queue

type Message struct {
	To      int
	Content string
}

type handler func(msg *Message)

type queue struct {
	ch chan *Message
}

func newQueue() *queue {
	return &queue{ch: make(chan *Message)}
}

type MessageQueue struct {
	queues map[string]*queue
}

func New() *MessageQueue {
	return &MessageQueue{
		queues: make(map[string]*queue),
	}
}

func (manager *MessageQueue) register(topic string) {
	if _, exists := manager.queues[topic]; !exists {
		manager.queues[topic] = newQueue()
	}
}

func (manager *MessageQueue) Subscribe(topic string, handler handler) {
	if _, exists := manager.queues[topic]; !exists {
		manager.register(topic)
	}
	go func() {
		for msg := range manager.queues[topic].ch {
			handler(msg)
		}
	}()
}

func (manager *MessageQueue) Publish(topic string, msg *Message) error {
	if _, exists := manager.queues[topic]; !exists {
		manager.register(topic)
	}
	manager.queues[topic].ch <- msg
	return nil
}
