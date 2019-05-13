package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pqnguyen/simple-chatapp/message"
	"log"
	"net"
)

type Server struct{}

func (server *Server) Start(port string) {
	fmt.Printf("Starting server at port %s \n", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error while starting server: %s", err)
	}
	defer listener.Close()
	manager := ClientManager{
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
