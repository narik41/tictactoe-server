package internal

import (
	"fmt"
	"log"

	"github.com/narik41/tictactoe-helper/core"
	"github.com/narik41/tictactoe-server/internal/decoder"
)

type MessageRouter struct {
	handlers map[core.Version1MessageType]MessageHandler
}
type HandlerResponse struct {
	MessageType core.Version1MessageType
	Payload     interface{}
	Broadcast   bool
	Recipients  []string
}

type MessageHandler interface {
	Handle(msg *decoder.DecodedMessage, sessionId string) (*HandlerResponse, error)
	RequiredStates() []SessionState
}

func NewMessageRouter() *MessageRouter {
	router := &MessageRouter{
		handlers: make(map[core.Version1MessageType]MessageHandler),
	}

	return router
}

func (r *MessageRouter) RegisterHandler(msgType core.Version1MessageType, handler MessageHandler) {
	r.handlers[msgType] = handler
}

func (r *MessageRouter) Route(msg *decoder.DecodedMessage, session *Session) (*HandlerResponse, error) {
	// Step 1: Get handler
	handler, exists := r.handlers[msg.MessageType]
	if !exists {
		return nil, fmt.Errorf("unknown message type: %s", msg.MessageType)
	}

	// Step 2: Validate session state
	if err := r.validateSessionState(handler, session); err != nil {
		return nil, err
	}

	// Step 3: Call handler
	log.Printf("Routing %s for session %s", msg.MessageType, session)
	response, err := handler.Handle(msg, session.Id)
	if err != nil {
		return nil, fmt.Errorf("handler failed: %w", err)
	}

	return response, nil
}

func (r *MessageRouter) validateSessionState(handler MessageHandler, session *Session) error {
	requiredStates := handler.RequiredStates()
	if requiredStates == nil {
		return nil // No state restriction
	}

	for _, allowedState := range requiredStates {
		if session.State == allowedState {
			return nil
		}
	}

	return fmt.Errorf("invalid session state: %s not allowed for this operation", session.State)
}
