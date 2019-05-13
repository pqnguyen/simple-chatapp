package main

import (
	"github.com/pqnguyen/simple-chatapp/backend/message_queue"
	"github.com/pqnguyen/simple-chatapp/backend/models"
	_ "github.com/pqnguyen/simple-chatapp/backend/models"
	"github.com/pqnguyen/simple-chatapp/backend/redis"
	"github.com/pqnguyen/simple-chatapp/backend/server"
	"github.com/pqnguyen/simple-chatapp/backend/session"
	"sync"
)

func main() {
	// redis server responsible for record which server keep connection
	redisSrv := redis.New()

	// message queue responsible for storing message
	messageQueue := message_queue.New()

	// session server responsible for routing message between server
	sessionConfig := &session.Config{
		Redis:        redisSrv,
		MessageQueue: messageQueue,
	}

	sessionSrv := session.New(sessionConfig)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		serverConfig := &server.Config{
			Session:      sessionSrv,
			Redis:        redisSrv,
			MessageQueue: messageQueue,
			Port:         ":8080",
		}
		srv := server.New(serverConfig)
		srv.Start()
		wg.Done()
	}()

	go func() {
		serverConfig := &server.Config{
			Session:      sessionSrv,
			Redis:        redisSrv,
			MessageQueue: messageQueue,
			Port:         ":8081",
		}
		srv := server.New(serverConfig)
		srv.Start()
		wg.Done()
	}()

	wg.Wait()
	defer models.DB.Close()
}
