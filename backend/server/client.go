package server

import (
	"bufio"
	"encoding/json"
	"github.com/pqnguyen/simple-chatapp/backend/models"
	"github.com/pqnguyen/simple-chatapp/message"
	"github.com/pqnguyen/simple-chatapp/types"
	"github.com/pqnguyen/simple-chatapp/utils"
	"log"
	"net"
)

type Client struct {
	clientManager *ClientManager
	uid           int
	conn          net.Conn
	data          chan message.Talk
}

func (client *Client) receive() {
	reader := bufio.NewReader(client.conn)
	for {
		buf, err := reader.ReadBytes('\n')
		if len(buf) == 0 {
			continue
		}
		if err != nil {
			log.Printf("error while read message: %s", err)
			continue
		}
		var msg message.Talk
		if err := json.Unmarshal(buf, &msg); err != nil {
			log.Printf("error while unmarshall message: %s", err)
			continue
		}
		client.handleTalk(&msg)
	}
}

func (client *Client) handleTalk(talk *message.Talk) {
	session := client.clientManager.server.config.Session
	receiver, ok := client.clientManager.getClient(talk.To)
	if !ok {
		session.Push(talk)
		return
	}
	receiver.data <- *talk
}

func (client *Client) send() {
	for {
		select {
		case msgTalk := <-client.data:
			if err := client.sendMessage(msgTalk); err != nil {
				client.clientManager.unregister <- *client
				models.SaveMessage(msgTalk)
			}
		}
	}
}

func (client *Client) sendPing() error {
	ping := message.NewPing()
	return sendMessage(client, ping)
}

func (client *Client) sendMessage(talk message.Talk) error {
	err := client.sendPing()
	if err == nil {
		return sendMessage(client, talk)
	}
	return err
}

func sendMessage(client *Client, msg interface{}) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error while ping to client: %s", err)
	}
	return sendRawMessage(client, buf)
}

func sendRawMessage(client *Client, buf []byte) error {
	buf = utils.WrapperMessage(buf)
	if _, err := client.conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func (client *Client) sendUnreadMessage() {
	var messages []models.Message
	models.DB.
		Where("receiver = ? and status = ?", client.uid, types.Unread).
		Order("id").Find(&messages)
	for _, msg := range messages {
		var talk message.Talk
		_ = json.Unmarshal([]byte(msg.Message), &talk)
		_ = client.sendMessage(talk)
		msg.Status = types.Read
		models.DB.Save(&msg)
	}
}
