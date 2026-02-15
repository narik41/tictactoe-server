package internal

import (
	"bufio"
	"io"
	"log"
)

type SessionState string

const (
	Guest          SessionState = "GUEST"
	LoggedIn       SessionState = "LOGGED_IN"
	WaitingForPair SessionState = "WAITING_FOR_PAIR"
	IN_GAME        SessionState = "IN_GAME"
)

type Session struct {
	Id           string `json:"id"` // Session id
	Client       *Client
	State        SessionState // session state
	Username     string       // user username
	CreatedAt    int64
	LastActivity int64
}

func (s *Session) ReadLoop(sessionManager *SessionManager, messageRouter *MessageRouter, rw *bufio.ReadWriter) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ReadLoop panicked for session %s: %v", s.Id, r)
		}
		log.Printf("ReadLoop exited for session %s", s.Id)
	}()
	decoder := NewMessageDecoder(rw)
	sender := NewResponseSender(sessionManager)
	for {
		log.Printf("Reading a message of client %s", s.Id)
		decodedMsg, err2 := decoder.Decode()
		if err2 != nil {
			if err2 == io.EOF {
				log.Printf("Session %s disconnected", s.Id)
				return
			}
			log.Printf("Decode error for session %s: %v", s.Id, err2)
			sender.SendError(s, "DECODE_ERROR", err2.Error())
			continue
		}
		log.Printf("Session %s received %s", s.Id, decodedMsg.MessageType)

		response, err2 := messageRouter.Route(decodedMsg, s)
		if err2 != nil {
			sender.SendError(s, "HANDLER_ERROR", err2.Error())
			continue
		}

		if response.Broadcast {
			sender.Broadcast(response.Recipients, response)
		} else {
			sender.Send(s, response)
		}

	}
}
