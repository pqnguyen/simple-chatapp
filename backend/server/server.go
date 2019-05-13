package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pqnguyen/simple-chatapp/backend/message_queue"
	"github.com/pqnguyen/simple-chatapp/backend/redis"
	"github.com/pqnguyen/simple-chatapp/backend/session"
	"github.com/pqnguyen/simple-chatapp/message"
	"log"
	"net"
)

type Config struct {
	Redis        *redis.Redis
	MessageQueue *message_queue.MessageQueue
	Session      *session.Session
	Port         string
}

type Server struct {
	config Config
}

func New(config *Config) *Server {
	return &Server{config: *config}
}

func (server *Server) Start() {
	port := server.config.Port
	fmt.Printf("Starting server at port %s \n", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error while starting server: %s", err)
	}
	defer listener.Close()
	manager := ClientManager{
		server:     server,
		clients:    make(map[int]Client),
		register:   make(chan Client),
		unregister: make(chan Client),
	}
	go manager.start()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error while accept connection: %s \n", err)
		}
		buf, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			log.Printf("error while read message: %s \n", err)
		}
		var register message.Register
		if err := json.Unmarshal(buf, &register); err != nil {
			log.Printf("error while unmarshal register message: %s", err)
		}
		// record which sever connect to client
		server.config.Redis.Put(register.UID, port)
		fmt.Printf("user %d connected \n", register.UID)
		client := Client{
			clientManager: &manager,
			uid:           register.UID,
			conn:          conn,
			data:          make(chan string),
		}
		manager.register <- client
		go client.receive()
		go client.send()
	}
}
