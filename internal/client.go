package internal

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/narik41/tictactoe-message/core"
)

type Client struct {
	ClientId string
	Conn     net.Conn // client connection

}

func NewClient(conn net.Conn, server *Server) *Client {
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

func (c *Client) ReadLoop() {
	rw := bufio.NewReadWriter(bufio.NewReader(c.Conn), bufio.NewWriter(c.Conn))
	for {
		log.Printf("Reading a message of client %s", c.ClientId)
		line, err := rw.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected: %v", c.ClientId, err)
			return
		}
		newLine := bytes.TrimSpace([]byte(line))
		fmt.Printf("=== RAW MESSAGE DEBUG ===\n")
		fmt.Printf("Length: %d bytes\n", len(line))
		fmt.Printf("Raw string: %s\n", string(line))
		fmt.Printf("First char: %q (byte: %d)\n", line[0], line[0])
		if len(line) > 1 {
			fmt.Printf("Second char: %q (byte: %d)\n", line[1], line[1])
		}
		fmt.Printf("========================\n")
		if len(line) > 0 && line[0] == '"' {
			// It's a JSON string, unmarshal to get the Base64 string
			var base64Str string

			if err := json.Unmarshal(newLine, &base64Str); err != nil {
				//return nil, &ProtocolError{
				//	Code:    "INVALID_JSON",
				//	Message: fmt.Sprintf("failed to parse outer JSON string: %v", err),
				//}
			}

			// Decode from Base64
			decodedBytes, _ := base64.StdEncoding.DecodeString(base64Str)

			newLine = decodedBytes
		}
		message, err := core.DecodeMessage(newLine)
		if err != nil {
			log.Printf("Invalid message from %s: %v", c.ClientId, err)
			continue
		}
		log.Printf("Client %s received message: %s", c.ClientId, message)
		//c.msgHandler.ProcessMessage(message)

		jsonBytes, err := json.Marshal(message.Payload)
		if err != nil {
			log.Println(err)
			return
		}

		var v1Msg *core.Version1MessagePayload
		if err := json.Unmarshal(jsonBytes, &v1Msg); err != nil {
			log.Println(err)
		}

		if v1Msg.MessageType == core.MSG_LOGIN_PAYLOAD {
			// return success response
			milli := time.Now().UnixMilli()
			loginReqPayload := &core.Version1MessagePayload{
				MessageType: core.MSG_LOGIN_RESPONSE,
				Payload: &core.Version1MessageLoginResponse{
					IsAuthenticated: true,
				},
			}

			ticTacToeMsg := core.TicTacToeMessage{
				MessageId: generateID(),
				Version:   "v1",
				Timestamp: milli,
				Payload:   loginReqPayload,
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

	}
}
