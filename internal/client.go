package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/narik41/tictactoe-message/core"
)

type Client struct {
	ClientId   string
	Conn       net.Conn // client connection
	msgHandler MessageHandler
}

func NewClient(conn net.Conn, server *Server, msgHandler MessageHandler) *Client {
	return &Client{
		ClientId:   generateID(),
		Conn:       conn,
		msgHandler: msgHandler,
	}
}

var idCounter int
var idMu sync.Mutex

func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), idCounter)
}

func (c *Client) ReadLoop() {
	rw := bufio.NewReadWriter(bufio.NewReader(c.Conn), bufio.NewWriter(c.Conn))
	for {

		line, err := rw.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected: %v", c.ClientId, err)
			return
		}

		message, err := core.DecodeMessage([]byte(line))
		if err != nil {
			log.Printf("Invalid message from %s: %v", c.ClientId, err)
			continue
		}
		log.Printf("Client %s received message: %s", c.ClientId, message)
		c.msgHandler.ProcessMessage(message)

	}
}
