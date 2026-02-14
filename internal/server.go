package internal

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/narik41/tictactoe-message/core"
)

type Server struct {
	listener  net.Listener
	clientMap map[string]*Client
	//msgHandler MessageHandler
}

func NewServer() *Server {
	return &Server{

		clientMap: make(map[string]*Client),
	}
}

func (s *Server) Start(addr string) error {
	log.Printf("Starting server on addr %s", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	log.Printf("Server started on %s", addr)

	log.Printf("Listening for connections on addr %s", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Client Accept error: %v", err)
			continue
		}
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	log.Printf("Handling connection from %s", conn.RemoteAddr())
	client := NewClient(conn, s)
	s.clientMap[client.ClientId] = client

	log.Printf("Client connected to %s", client.ClientId)

	// ask username and password
	milli := time.Now().UnixMilli()
	loginReqPayload := &core.Version1MessagePayload{
		MessageType: core.MSG_LOGIN_REQUEST,
	}

	ticTacToeMsg := core.TicTacToeMessage{
		MessageId: generateID(),
		Version:   "v1",
		Timestamp: milli,
		Payload:   loginReqPayload,
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	if err := json.NewEncoder(rw.Writer).Encode(ticTacToeMsg); err != nil {
		log.Printf("Encoding error: %v", err)
		return
	}

	err := rw.Flush()
	if err != nil {
		log.Printf("Flush error: %v", err)
		return
	}

	go client.ReadLoop()
}
