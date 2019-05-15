package session

import (
	"github.com/pqnguyen/simple-chatapp/backend/message_queue"
	"github.com/pqnguyen/simple-chatapp/backend/models"
	"github.com/pqnguyen/simple-chatapp/backend/redis"
	"github.com/pqnguyen/simple-chatapp/message"
	"log"
)

type Config struct {
	Redis        *redis.Redis
	MessageQueue *message_queue.MessageQueue
}

type Session struct {
	config Config
}

func New(config *Config) *Session {
	return &Session{
		config: *config,
	}
}

func (session *Session) Push(talk *message.Talk) {
	redis := session.config.Redis
	messageQueue := session.config.MessageQueue
	topic := redis.Get(talk.To)
	if topic == "" {
		models.SaveMessage(*talk)
		return
	}
	if err := messageQueue.Publish(topic, talk); err != nil {
		log.Printf("error while publish message to queue")
	}
}
