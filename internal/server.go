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
	listener           net.Listener
	sessionManager     *SessionManager
	gameSessionManager *GameSessionManager
	clientMap          map[string]*Client
	waitingClient      []*Client
	sessionMap         map[string]*Session
	//msgHandler MessageHandler
}

func NewServer() *Server {
	return &Server{
		clientMap:      make(map[string]*Client),
		sessionMap:     make(map[string]*Session),
		waitingClient:  make([]*Client, 0),
		sessionManager: NewSessionManager(),
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
	client := NewClient(conn)
	session := s.sessionManager.CreateSession(client)
	defer s.sessionManager.RemoveSession(session.Id)

	log.Printf("Session %s created for client", session.Id)

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

	go session.ReadLoop(s.sessionManager)
}

func (s *Server) processClient(newClient *Client) {
	log.Printf("Finding opponent for the client %s", newClient)

	// check for waiting list
	if len(s.waitingClient) > 0 {
		opponent := s.waitingClient[0]
		sessionId := "1"
		gameSession := NewGameSession(sessionId)
		gameSession.AddPlayer(sessionId, newClient.ClientId, newClient)
		gameSession.AddPlayer(sessionId, opponent.ClientId, opponent)
		gameSession.Start()
		s.waitingClient = s.waitingClient[1:]

		// relay the message to both client that game session is started
	} else {
		s.waitingClient = append(s.waitingClient, newClient)
	}

}
