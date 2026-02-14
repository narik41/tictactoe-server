package internal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/narik41/tictactoe-message/core"
	"github.com/narik41/tictactoe-server/internal/repo"
)

type AuthService struct {
	userRepo repo.UserRepo
}

func NewAuthService(userRepo repo.UserRepo) AuthService {
	return AuthService{
		userRepo: userRepo,
	}
}

func (a AuthService) Handle(msg *DecodedMessage, session *Session) (*HandlerResponse, error) {
	log.Println("Handling the auth request.")
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	// Unmarshal to v1 payload structure
	var loginPayload core.Version1MessageLoginPayload
	if err := json.Unmarshal(jsonBytes, &loginPayload); err != nil {
		return nil, err
	}

	isExists := a.userRepo.GetByUsername(loginPayload.Username)
	if !isExists {
		return nil, fmt.Errorf("user not found")
	}
	session.State = LoggedIn
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

func (a AuthService) RequiredStates() []SessionState {
	return []SessionState{
		Guest,
	}
}
