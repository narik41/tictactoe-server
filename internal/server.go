package internal

import (
	"log"
	"net"
)

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	return &Server{}
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

}
