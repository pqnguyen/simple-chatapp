package server

type ClientManager struct {
	clients    map[int]Client
	register   chan Client
	unregister chan Client
}

func (manager *ClientManager) registerClient(client Client) {
	manager.clients[client.uid] = client
}

func (manager *ClientManager) unregisterClient(uid int) {
	if _, ok := manager.clients[uid]; ok {
		delete(manager.clients, uid)
	}
}

func (manager *ClientManager) getClient(uid int) (Client, bool) {
	client, ok := manager.clients[uid]
	return client, ok
}

func (manager *ClientManager) start() {
	for {
		select {
		case client := <-manager.register:
			manager.registerClient(client)
		case client := <-manager.unregister:
			manager.unregisterClient(client.uid)
		}
	}
}
