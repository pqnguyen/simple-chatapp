package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/pqnguyen/simple-chatapp/message"
	"log"
	"net"
	"os"
	"sync"
)

var (
	identity int
	receiver int
	port     string
)

type Client struct {
	uid  int
	conn net.Conn
}

func (client *Client) send(data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return errors.New(fmt.Sprintf("error while marshal message: %s", err))
	}
	if ok := bytes.HasSuffix(buf, []byte{'\n'}); !ok {
		buf = append(buf, '\n')
	}
	if _, err := client.conn.Write(buf); err != nil {
		return errors.New(fmt.Sprintf("error while write message to buffer: %s", err))
	}
	fmt.Print(string(buf))
	return nil
}

func (client *Client) sendRegister() {
	register := message.NewRegister(client.uid)
	if err := client.send(register); err != nil {
		log.Fatalf("error while send register: %s", err)
	}
}

func (client *Client) sendTalk(to int, content string) {
	if to == client.uid {
		log.Printf("can't send to yourself")
	}
	talk := message.NewTalk(client.uid, to, content)
	if err := client.send(talk); err != nil {
		log.Printf("error while send talk: %s", err)
	}
}

func (client *Client) sendUnregister() {
	unregister := message.NewUnregister(client.uid)
	if err := client.send(unregister); err != nil {
		log.Printf("error while send unregister: %s", err)
	}
}

func (client *Client) listen() {
	reader := bufio.NewReader(client.conn)
	for {
		//buf, err := reader.ReadBytes('\n')
		buf := make([]byte, 1024)
		_, err := reader.Read(buf)
		if err != nil {
			log.Fatalf("error while read message from connection: %s", err)
		}
		if len(buf) == 0 {
			continue
		}
		fmt.Printf("Received message: %s", string(buf))
	}
}

func (client *Client) start() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content := scanner.Text()
		client.sendTalk(receiver, content)
	}
}

func NewClient(uid int) *Client {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalf("error while connect to server on port 8080: %s", err)
	}
	client := Client{
		uid:  uid,
		conn: conn,
	}
	client.sendRegister()
	return &client
}

func main() {
	flag.IntVar(&identity, "id", 0, "user identity")
	flag.IntVar(&receiver, "receiver", 0, "receiver identity")
	flag.StringVar(&port, "port", ":8080", "server you want to connect to")
	flag.Parse()

	if identity == 0 || receiver == 0 {
		log.Fatal("missing args")
	}
	fmt.Printf("user %d: Chatting with %d \n", identity, receiver)
	client := NewClient(identity)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		client.listen()
		wg.Done()
	}()
	client.start()
	wg.Wait()
}
