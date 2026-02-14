package internal

import (
	"log"

	"github.com/narik41/tictactoe-message/core"
)

type MessageHandler interface {
	ProcessMessage(msg *core.TicTacToeMessage) error
}

type MessageHandlerImpl struct {
	v1MsgHandler Version1MsgHandler
}

func NewMessageHandler(v1MsgHandler Version1MsgHandler) MessageHandler {
	return MessageHandlerImpl{
		v1MsgHandler: v1MsgHandler,
	}
}

func (m MessageHandlerImpl) ProcessMessage(msg *core.TicTacToeMessage) error {
	switch msg.Version {
	case "v1.0.0":
		return m.v1MsgHandler.Version1MessageHandler(msg.Payload)
	default:
		log.Printf("Received: %s", msg.MessageId)
	}
	return nil
}
