package internal

import (
	"log"

	"github.com/narik41/tictactoe-helper/core"
)

type HeartbeatHandler struct{}

func NewHeartbeatHandler() *HeartbeatHandler {
	return &HeartbeatHandler{}
}

func (h *HeartbeatHandler) RequiredStates() []SessionState {
	return nil
}

func (h *HeartbeatHandler) Handle(msg *DecodedMessage, session *Session) (*HandlerResponse, error) {
	log.Printf("Heartbeat from session %s", session.Id)

	session.LastActivity = core.GetNPTToUtcInMillisecond()
	return &HandlerResponse{
		MessageType: "HEARTBEAT_RESPONSE",
		Payload:     nil,
	}, nil
}
