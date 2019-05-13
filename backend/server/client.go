package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pqnguyen/simple-chatapp/backend/models"
	"github.com/pqnguyen/simple-chatapp/message"
	"github.com/pqnguyen/simple-chatapp/types"
	"log"
	"net"
)

type Client struct {
	clientManager *ClientManager
	uid           int
	conn          net.Conn
	data          chan string
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
		session.Push(talk.To, talk.Content)
		return
	}
	receiver.data <- talk.Content
}

func (client *Client) send() {
	for {
		select {
		case content := <-client.data:
			if err := client.sendMessage(content); err != nil {
				client.clientManager.unregister <- *client
				models.DB.Create(&models.Message{
					Message:  content,
					Receiver: client.uid,
					Sender:   0,
				})
			}
		}
	}
}

func (client *Client) sendMessage(content string) error {
	buf := []byte(content)
	if ok := bytes.HasSuffix(buf, []byte{'\n'}); !ok {
		buf = append(buf, '\n')
	}
	if _, err := client.conn.Write(buf); err != nil {
		client.clientManager.unregister <- *client
		models.DB.Create(&models.Message{
			Message:  content,
			Receiver: client.uid,
			Sender:   0,
		})
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
		_ = client.sendMessage(msg.Message)
		msg.Status = types.Read
		models.DB.Save(&msg)
	}
}
