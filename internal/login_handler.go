package internal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/narik41/tictactoe-helper/core"
	"github.com/narik41/tictactoe-server/internal/decoder"
	"github.com/narik41/tictactoe-server/internal/repo"
)

type LoginHandler struct {
	userRepo       repo.UserRepo
	queue          *SessionQueue
	sessionManager *SessionManager
}

func NewLoginHandler(userRepo repo.UserRepo, queue *SessionQueue, sessionManager *SessionManager) LoginHandler {
	return LoginHandler{
		userRepo:       userRepo,
		queue:          queue,
		sessionManager: sessionManager,
	}
}

func (a LoginHandler) Handle(msg *decoder.DecodedMessage, sessionId string) (*HandlerResponse, error) {
	log.Println("Handling the auth request.")
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var loginPayload core.Version1MessageLoginPayload
	if err := json.Unmarshal(jsonBytes, &loginPayload); err != nil {
		return nil, err
	}

	isExists := a.userRepo.GetByUsername(loginPayload.Username)
	if !isExists {
		return nil, fmt.Errorf("user not found")
	}

	clientSession, _ := a.sessionManager.GetSession(sessionId)

	// add session to queue
	a.queue.Enqueue(clientSession)
	return &HandlerResponse{
		MessageType: core.MSG_LOGIN_RESPONSE,
		Payload: &core.Version1MessageLoginResponse{
			IsAuthenticated: true,
			Message:         loginPayload.Username,
			PlayerId:        loginPayload.Username,
		},
		Broadcast: false,
	}, nil
}

func (a LoginHandler) RequiredStates() []SessionState {
	return []SessionState{
		Guest,
	}
}
