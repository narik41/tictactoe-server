package internal

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"github.com/narik41/tictactoe-message/core"
)

type SessionState string

const (
	Guest    SessionState = "GUEST"
	LoggedIn SessionState = "LOGGED_IN"
)

type Session struct {
	Id        string `json:"id"` // Session id
	Client    *Client
	State     SessionState // session state
	Username  string       // user username
	CreatedAt int64
}

func (s *Session) ReadLoop(sessionManager *SessionManager) {

	decoder := NewMessageDecoder(s.Client.Conn)
	messageRouter := NewMessageRouter()
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

		//line, err := rw.ReadString('\n')
		//if err != nil {
		//	if err == io.EOF {
		//		log.Printf("Client %s disconnected: %v", c.Id, err)
		//		return
		//	}
		//	log.Printf("Error %v", err)
		//	return
		//}
		//newLine := bytes.TrimSpace([]byte(line))
		//if len(line) > 0 && line[0] == '"' {
		//	var base64Str string
		//
		//	if err := json.Unmarshal(newLine, &base64Str); err != nil {
		//		//return nil, &ProtocolError{
		//		//	Code:    "INVALID_JSON",
		//		//	Message: fmt.Sprintf("failed to parse outer JSON string: %v", err),
		//		//}
		//	}
		//
		//	// Decode from Base64
		//	decodedBytes, _ := base64.StdEncoding.DecodeString(base64Str)
		//
		//	newLine = decodedBytes
		//}
		//message, err := core.DecodeMessage(newLine)
		//if err != nil {
		//	log.Printf("Invalid message from %s: %v", c.Id, err)
		//	continue
		//}
		//log.Printf("Client %s received message: %s", c.Id, message)
		////c.msgHandler.ProcessMessage(message)
		//
		//jsonBytes, err := json.Marshal(message.Payload)
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
		//
		//var v1Msg *core.Version1MessagePayload
		//if err := json.Unmarshal(jsonBytes, &v1Msg); err != nil {
		//	log.Println(err)
		//}
		//
		//if v1Msg.MessageType == core.MSG_LOGIN_PAYLOAD {
		//
		//	milli := time.Now().UnixMilli()
		//	loginReqPayload := &core.Version1MessagePayload{
		//		MessageType: core.MSG_LOGIN_RESPONSE,
		//		Payload: &core.Version1MessageLoginResponse{
		//			IsAuthenticated: true,
		//		},
		//	}
		//
		//	ticTacToeMsg := core.TicTacToeMessage{
		//		MessageId: generateID(),
		//		Version:   "v1",
		//		Timestamp: milli,
		//		Payload:   loginReqPayload,
		//	}
		//
		//	if err := json.NewEncoder(rw.Writer).Encode(ticTacToeMsg); err != nil {
		//		log.Printf("Encoding error: %v", err)
		//		return
		//	}
		//
		//	err := rw.Flush()
		//	if err != nil {
		//		log.Printf("Flush error: %v", err)
		//		return
		//	}
		//
		//} else if v1Msg.MessageType == core.PLAYER_MOVE {
		//
		//}

	}
}

func (c *Session) RelayGameStarted(playerInfo *PlayerInfo) {
	rw := bufio.NewReadWriter(bufio.NewReader(c.Client.Conn), bufio.NewWriter(c.Client.Conn))

	gameStartPayload := &core.Version1MessagePayload{
		MessageType: core.GAME_START,
		Payload: &core.Version1GameStartPayload{
			YourTurn:   playerInfo.MyTurn,
			YourSymbol: string(playerInfo.Symbol),
		},
	}

	ticTacToeMsg := core.TicTacToeMessage{
		MessageId: generateID(),
		Version:   "v1",
		Payload:   gameStartPayload,
	}

	if err := json.NewEncoder(rw.Writer).Encode(ticTacToeMsg); err != nil {
		log.Printf("Encoding error: %v", err)
		return
	}

	err := rw.Flush()
	if err != nil {
		log.Printf("Flush error: %v", err)
		return
	}

}
