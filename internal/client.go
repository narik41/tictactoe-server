package internal

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	ClientId string
	Conn     net.Conn // client connection
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		ClientId: generateID(),
		Conn:     conn,
	}
}

var idCounter int
var idMu sync.Mutex

func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), idCounter)
}

//
//func (c *Client) ReadLoop() {
//	rw := bufio.NewReadWriter(bufio.NewReader(c.Conn), bufio.NewWriter(c.Conn))
//	for {
//		log.Printf("Reading a message of client %s", c.ClientId)
//		line, err := rw.ReadString('\n')
//		if err != nil {
//			if err == io.EOF {
//				log.Printf("Client %s disconnected: %v", c.ClientId, err)
//				return
//			}
//			log.Printf("Error %v", err)
//			return
//		}
//		newLine := bytes.TrimSpace([]byte(line))
//		if len(line) > 0 && line[0] == '"' {
//			var base64Str string
//
//			if err := json.Unmarshal(newLine, &base64Str); err != nil {
//				//return nil, &ProtocolError{
//				//	Code:    "INVALID_JSON",
//				//	Message: fmt.Sprintf("failed to parse outer JSON string: %v", err),
//				//}
//			}
//
//			// Decode from Base64
//			decodedBytes, _ := base64.StdEncoding.DecodeString(base64Str)
//
//			newLine = decodedBytes
//		}
//		message, err := core.DecodeMessage(newLine)
//		if err != nil {
//			log.Printf("Invalid message from %s: %v", c.ClientId, err)
//			continue
//		}
//		log.Printf("Client %s received message: %s", c.ClientId, message)
//		//c.msgHandler.ProcessMessage(message)
//
//		jsonBytes, err := json.Marshal(message.Payload)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//
//		var v1Msg *core.Version1MessagePayload
//		if err := json.Unmarshal(jsonBytes, &v1Msg); err != nil {
//			log.Println(err)
//		}
//
//		if v1Msg.MessageType == core.MSG_LOGIN_PAYLOAD {
//
//			milli := time.Now().UnixMilli()
//			loginReqPayload := &core.Version1MessagePayload{
//				MessageType: core.MSG_LOGIN_RESPONSE,
//				Payload: &core.Version1MessageLoginResponse{
//					IsAuthenticated: true,
//				},
//			}
//
//			ticTacToeMsg := core.TicTacToeMessage{
//				MessageId: generateID(),
//				Version:   "v1",
//				Timestamp: milli,
//				Payload:   loginReqPayload,
//			}
//
//			if err := json.NewEncoder(rw.Writer).Encode(ticTacToeMsg); err != nil {
//				log.Printf("Encoding error: %v", err)
//				return
//			}
//
//			err := rw.Flush()
//			if err != nil {
//				log.Printf("Flush error: %v", err)
//				return
//			}
//			//c.server.processClient(c)
//		} else if v1Msg.MessageType == core.PLAYER_MOVE {
//
//		}
//
//	}
//}
//
//func (c *Client) RelayGameStarted(playerInfo *PlayerInfo) {
//	rw := bufio.NewReadWriter(bufio.NewReader(c.Conn), bufio.NewWriter(c.Conn))
//
//	gameStartPayload := &core.Version1MessagePayload{
//		MessageType: core.GAME_START,
//		Payload: &core.Version1GameStartPayload{
//			YourTurn:   playerInfo.MyTurn,
//			YourSymbol: string(playerInfo.Symbol),
//		},
//	}
//
//	ticTacToeMsg := core.TicTacToeMessage{
//		MessageId: generateID(),
//		Version:   "v1",
//		Payload:   gameStartPayload,
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
//}
