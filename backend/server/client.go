package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pqnguyen/simple-chatapp/message"
	"log"
	"net"
)

type Client struct {
	uid           int
	conn          net.Conn
	clientManager *ClientManager
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
		}
		var talk message.Talk
		if err := json.Unmarshal(buf, &talk); err != nil {
			log.Printf("error while unmarshall message: %s", err)
		}
		fmt.Printf("Got message from %d: %s \n", client.uid, talk.Content)
		client, ok := client.clientManager.getClient(talk.To)
		if !ok {
			log.Printf("the receiver is not exists")
			continue
		}
		client.data <- talk.Content
	}
}

func (client *Client) send() {
	for {
		select {
		case content := <-client.data:
			buf := []byte(content)
			if ok := bytes.HasSuffix(buf, []byte{'\n'}); !ok {
				buf = append(buf, '\n')
			}
			if _, err := client.conn.Write(buf); err != nil {
				log.Printf("error while send message: %s", err)
			}
		}
	}
}
