package client

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"os"
)

type Message struct {
	Sender string
	Data   string
}

type Client struct {
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	encoder *gob.Encoder
	decoder *gob.Decoder
	name    string
	status  int
}

func (client *Client) Go() {
	go client.ListenRoom()

	// forever - read stdin and send it
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		client.SendText(text)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("reading standard input:", err)
	}
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) SendText(text string) {
	err := client.encoder.Encode(Message{
		Sender: client.name,
		Data:   text,
	})
	client.writer.Flush()

	if err != nil {
		log.Fatalf("encoder error: %s", err)
	}

	log.Printf("<<\t\t\t\t '%s'", text)
}

func (client *Client) ListenRoom() {
	var m Message
	for client.status == 1 {
		err := client.decoder.Decode(&m)
		if err != nil {
			log.Fatalf("decoder error: %s", err)
		}
		if m.Sender != client.name {
			log.Printf(">> '(@%s) %s'", m.Sender, m.Data)
			m.Data = ""
		}
	}
}

func NewClient(addr string, name string) *Client {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[client] ")
	log.Println("started")

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	log.Println("connected")

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	encoder := gob.NewEncoder(writer)
	decoder := gob.NewDecoder(reader)

	client := &Client{
		conn:    conn,
		reader:  reader,
		writer:  writer,
		encoder: encoder,
		decoder: decoder,
		name:    name,
		status:  1,
	}
	return client
}

/*
func main() {
	client := NewClient(":6666", os.Args[1])
	defer client.Close()

	client.Go()
}
*/
