package server

import (
	"github.com/pqnguyen/simple-chatapp/backend/models"
	"github.com/pqnguyen/simple-chatapp/message"
	"log"
)

type ClientManager struct {
	server     *Server
	clients    map[int]Client
	register   chan Client
	unregister chan Client
}

func (manager *ClientManager) registerClient(client Client) {
	manager.clients[client.uid] = client
	client.sendUnreadMessage()
}

func (manager *ClientManager) unregisterClient(uid int) {
	if _, ok := manager.clients[uid]; ok {
		_ = manager.clients[uid].conn.Close()
		delete(manager.clients, uid)
	}
}

func (manager *ClientManager) getClient(uid int) (Client, bool) {
	client, ok := manager.clients[uid]
	return client, ok
}

func (manager *ClientManager) start() {
	serverConfig := manager.server.config
	serverConfig.MessageQueue.Subscribe(serverConfig.Port, manager.forward)
	for {
		select {
		case client := <-manager.register:
			manager.registerClient(client)
		case client := <-manager.unregister:
			manager.unregisterClient(client.uid)
		}
	}
}

func (manager *ClientManager) forward(message *message.Talk) {
	client, ok := manager.getClient(message.To)
	if !ok {
		models.SaveMessage(*message)
		log.Printf("user %d doesn't exists", message.To)
		return
	}
	if err := client.sendMessage(*message); err != nil {
		client.clientManager.unregister <- client
		models.SaveMessage(*message)
		log.Printf("error while send message to client: %s", err)
	}
}
