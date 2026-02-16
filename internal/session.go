package internal

import (
	"bufio"
	"io"
	"log"

	"github.com/narik41/tictactoe-helper/core"
	"github.com/narik41/tictactoe-server/internal/decoder"
)

type SessionState string

const (
	Guest          SessionState = "GUEST"
	LoggedIn       SessionState = "LOGGED_IN"
	WaitingForPair SessionState = "WAITING_FOR_PAIR"
	IN_GAME        SessionState = "IN_GAME"
)

type Session struct {
	Id           string `json:"id"`
	Client       *Client
	State        SessionState
	Username     string
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
	msgDecoder := decoder.NewMessageDecoder(rw)
	msgSender := NewResponseSender(sessionManager)
	for {
		log.Printf("Reading a message of client %s", s.Id)
		decodedMsg, err2 := msgDecoder.Decode()
		if err2 != nil {
			if err2 == io.EOF {
				log.Printf("Session %s disconnected", s.Id)
				return
			}
			log.Printf("Decode error for session %s: %v", s.Id, err2)
			msgSender.SendError(s, "DECODE_ERROR", err2.Error())
			continue
		}
		log.Printf("Session %s received %s", s.Id, decodedMsg.MessageType)

		response, err2 := messageRouter.Route(decodedMsg, s)
		if err2 != nil {
			msgSender.SendError(s, "HANDLER_ERROR", err2.Error())
			continue
		}

		// based on response update the session state
		if response != nil {
			if response.MessageType == core.MSG_LOGIN_RESPONSE {
				s.State = LoggedIn
			}
		}

		if response.Broadcast {
			msgSender.Broadcast(response.Recipients, response)
		} else {
			msgSender.Send(s, response)
		}

	}
}
