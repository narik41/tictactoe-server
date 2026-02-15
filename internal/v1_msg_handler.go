package internal

import (
	"errors"

	"github.com/narik41/tictactoe-message/core"
)

type Version1MsgHandler struct{}

func NewVersion1MsgHandler() Version1MsgHandler {
	return Version1MsgHandler{}
}

func (m *Version1MsgHandler) Version1MessageHandler(v1MsgPayload interface{}) error {
	v1Msg, ok := v1MsgPayload.(core.Version1MessagePayload)
	if !ok {
		return errors.New("v1 payload is not a TicTacToeMessage")
	}

	switch v1Msg.MessageType {
	case core.MSG_LOGIN_REQUEST:
		//return m.doAuth()
	default:
		return errors.New("v1 payload is not a TicTacToeMessage")
	}
	return nil
}

//if message.Version == "v1" {
//
//	payload := message.Payload.(core.Version1MessagePayload)
//	if payload.MessageType == core.MSG_LOGIN_PAYLOAD {
//		msg := core.TicTacToeMessage{
//			Version: message.Version,
//			Payload: core.Version1MessagePayload{
//				MessageType: core.MSG_LOGIN_RESPONSE,
//				Payload: &core.Version1MessageLoginResponse{
//					IsAuthenticated: true,
//					Message:         "Successfully authenticated",
//					PlayerId:        c.ClientId,
//				},
//			},
//		}
//
//		if err := json.NewEncoder(rw.Writer).Encode(msg); err != nil {
//			log.Printf("Error encoding message: %v", err)
//		}
//
//		err := rw.Flush()
//		if err != nil {
//			log.Printf("Error flushing writer: %v", err)
//			return
//		}
//	}
//}
