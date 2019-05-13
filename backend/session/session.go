package session

import (
	"github.com/pqnguyen/simple-chatapp/backend/message_queue"
	"github.com/pqnguyen/simple-chatapp/backend/models"
	"github.com/pqnguyen/simple-chatapp/backend/redis"
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

func (session *Session) Push(to int, msg string) {
	redis := session.config.Redis
	messageQueue := session.config.MessageQueue
	topic := redis.Get(to)
	if topic == "" {
		models.DB.Create(&models.Message{
			Receiver: to,
			Sender:   0,
			Message:  msg,
		})
		return
	}
	message := message_queue.Message{To: to, Content: msg}
	if err := messageQueue.Publish(topic, &message); err != nil {
		log.Printf("error while publish message to queue")
	}
}
