package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/narik41/tictactoe-helper/core"
)

type ResponseSender struct {
	sessionManager *SessionManager
}

func NewResponseSender(sessionManager *SessionManager) *ResponseSender {
	return &ResponseSender{
		sessionManager: sessionManager,
	}
}

func (rs *ResponseSender) Send(session *Session, response *HandlerResponse) error {

	msgBytes, err := rs.encodeMessage(response.MessageType, response.Payload)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}

	if err := rs.writeToConn(session.Client.Conn, msgBytes); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	log.Printf("Sent %s to session %s", response.MessageType, session.Id)
	return nil
}

func (rs *ResponseSender) Broadcast(recipientIDs []string, response *HandlerResponse) error {

	msgBytes, err := rs.encodeMessage(response.MessageType, response.Payload)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}

	// Send to each recipient
	var sendErrors []error
	for _, sessionID := range recipientIDs {
		session, exists := rs.sessionManager.GetSession(sessionID)
		if !exists {
			log.Printf("Session %s not found for broadcast", sessionID)
			continue
		}

		if err := rs.writeToConn(session.Client.Conn, msgBytes); err != nil {
			log.Printf("Failed to broadcast to session %s: %v", sessionID, err)
			sendErrors = append(sendErrors, err)
			continue
		}

		log.Printf("Broadcasted %s to session %s", response.MessageType, sessionID)
	}

	if len(sendErrors) > 0 {
		return fmt.Errorf("failed to send to %d recipients", len(sendErrors))
	}

	return nil
}

func (rs *ResponseSender) SendError(session *Session, errorCode, errorMessage string) error {
	errorPayload := &core.Version1MessagePayload{
		MessageType: core.ERROR,
		Payload: map[string]interface{}{
			"code":    errorCode,
			"message": errorMessage,
		},
	}

	msgBytes, err := rs.encodeMessage("ERROR", errorPayload)
	if err != nil {
		return fmt.Errorf("failed to encode error: %w", err)
	}

	if err := rs.writeToConn(session.Client.Conn, msgBytes); err != nil {
		return fmt.Errorf("failed to send error: %w", err)
	}

	log.Printf("Sent error to session %s: %s", session.Id, errorCode)
	return nil
}

func (rs *ResponseSender) encodeMessage(messageType core.Version1MessageType, payload interface{}) ([]byte, error) {

	v1Payload := &core.Version1MessagePayload{
		MessageType: messageType,
		Payload:     payload,
	}

	msg := core.TicTacToeMessage{
		MessageId: core.UUID("msg"),
		Version:   "v1",
		Timestamp: time.Now().UnixMilli(),
		Payload:   v1Payload,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	msgBytes = append(msgBytes, '\n')

	return msgBytes, nil
}

func (rs *ResponseSender) writeToConn(conn net.Conn, data []byte) error {
	writer := bufio.NewWriter(conn)

	n, err := writer.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf("incomplete write: wrote %d of %d bytes", n, len(data))
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
