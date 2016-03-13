package server

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
)

// Message definition
// Type: 0(msg) 1(assign client id)
type Message struct {
	//	Type   int
	Sender string
	Data   string
}

type Client struct {
	chatRoom *ChatRoom
	incoming chan *Message
	outgoing chan *Message
	reader   *bufio.Reader
	writer   *bufio.Writer
	encoder  *gob.Encoder
	decoder  *gob.Decoder
	name     string
	status   int
}

// ListenSocket handle input from socket and put it into channel
func (client *Client) ListenSocket() {
	for client.status == 1 {
		var m Message
		err := client.decoder.Decode(&m)
		if err != nil {
			log.Printf("decoder error: %s", err)
			client.chatRoom.DeactivateClient(client)
		}

		client.incoming <- &m
		log.Printf("<< '(@%s) %s'", m.Sender, m.Data)
	}
}

func (client *Client) ListenChannel() {
	for {
		select {
		case m := <-client.outgoing:
			err := client.encoder.Encode(*m)
			client.writer.Flush()
			if err != nil {
				log.Fatalf("encoder error: %s", err)
			}

		case m := <-client.incoming:
			client.chatRoom.incoming <- m
		}
	}
}

func (client *Client) Go() {
	go client.ListenSocket()
	go client.ListenChannel()
}

func NewClient(connection net.Conn, chatRoom *ChatRoom) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)
	encoder := gob.NewEncoder(writer)
	decoder := gob.NewDecoder(reader)

	client := &Client{
		incoming: make(chan *Message),
		outgoing: make(chan *Message),
		reader:   reader,
		writer:   writer,
		encoder:  encoder,
		decoder:  decoder,
		chatRoom: chatRoom,
		status:   1,
	}
	return client
}

type ChatRoom struct {
	clients  []*Client
	joins    chan net.Conn
	incoming chan *Message
	outgoing chan *Message
	listener net.Listener
}

func (chatRoom *ChatRoom) Broadcast(msg *Message) {
	log.Printf("broadcasting: (@%s) %s", msg.Sender, msg.Data)

	for _, client := range chatRoom.clients {
		if client.status == 1 {
			client.outgoing <- msg
		}
	}
}

func (chatRoom *ChatRoom) Join(connection net.Conn) {
	client := NewClient(connection, chatRoom)
	client.Go()

	chatRoom.clients = append(chatRoom.clients, client)

	log.Printf("new client join")
}

func (chatRoom *ChatRoom) ListenChannel() {
	for {
		select {
		case msg := <-chatRoom.incoming:
			chatRoom.Broadcast(msg)

		case conn := <-chatRoom.joins:
			chatRoom.Join(conn)
		}
	}
}

func (chatRoom *ChatRoom) DeactivateClient(client *Client) {
	var i, length int
	length = len(chatRoom.clients)

	for i = 0; i < length; i++ {
		if chatRoom.clients[i] == client {
			chatRoom.clients[i].status = 0
		}
	}
}

func NewChatRoom(addr string) *ChatRoom {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[server] ")
	log.Println("started")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net.Listen error: %s", err.Error())
	}

	chatRoom := &ChatRoom{
		clients:  make([]*Client, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan *Message),
		outgoing: make(chan *Message),
		listener: listener,
	}

	return chatRoom
}

func (chatRoom *ChatRoom) Go() {
	go chatRoom.ListenChannel()

	for {
		conn, _ := chatRoom.listener.Accept()
		log.Println("new connection")
		chatRoom.joins <- conn
	}
}

/*
func main() {
	chatRoom := NewChatRoom(":6666")
	chatRoom.Go()
}
*/
