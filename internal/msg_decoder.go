package internal

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/narik41/tictactoe-message/core"
)

type MessageDecoder struct {
	reader *bufio.Reader
}

func NewMessageDecoder(conn io.Reader) *MessageDecoder {
	return &MessageDecoder{
		reader: bufio.NewReader(conn),
	}
}

func (d *MessageDecoder) Decode() (*DecodedMessage, error) {
	// Step 1: Read line until newline
	line, err := d.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	// Step 2: Trim whitespace
	line = bytes.TrimSpace(line)

	if len(line) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	// Step 3: Check if base64 encoded (wrapped in quotes)
	if line[0] == '"' {
		var base64Str string
		if err := json.Unmarshal(line, &base64Str); err != nil {
			return nil, fmt.Errorf("failed to parse base64 wrapper: %w", err)
		}

		// Decode from base64
		decoded, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}

		line = decoded
	}

	// Step 4: Parse base message structure
	var msg core.TicTacToeMessage
	if err := json.Unmarshal(line, &msg); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Step 5: Validate required fields
	//if msg.MessageId == "" {
	//	return nil, fmt.Errorf("message_id is required")
	//}
	//if msg.Version == "" {
	//	return nil, fmt.Errorf("version is required")
	//}
	//if msg.Timestamp == 0 {
	//	return nil, fmt.Errorf("timestamp is required")
	//}

	// Step 6: Extract version-specific payload
	var messageType core.Version1MessageType
	var payloadData interface{}

	switch msg.Version {
	case "v1":
		v1Payload, err := d.decodeV1Payload(msg.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode v1 payload: %w", err)
		}
		messageType = v1Payload.MessageType
		payloadData = v1Payload.Payload

	default:
		return nil, fmt.Errorf("unsupported version: %s", msg.Version)
	}

	// Step 7: Return decoded message
	return &DecodedMessage{
		MessageId:   msg.MessageId,
		Version:     msg.Version,
		MessageType: messageType,
		Payload:     payloadData,
		Timestamp:   msg.Timestamp,
	}, nil
}

// decodeV1Payload decodes version 1 payload
func (d *MessageDecoder) decodeV1Payload(payload interface{}) (*core.Version1MessagePayload, error) {
	// Convert to JSON bytes
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Unmarshal to v1 payload structure
	var v1Payload core.Version1MessagePayload
	if err := json.Unmarshal(jsonBytes, &v1Payload); err != nil {
		return nil, err
	}

	// Validate message type
	if v1Payload.MessageType == "" {
		return nil, fmt.Errorf("message_type is required")
	}

	return &v1Payload, nil
}

// DecodedMessage represents a fully decoded message
type DecodedMessage struct {
	MessageId   string
	Version     string
	MessageType core.Version1MessageType // e.g., "LOGIN_REQUEST", "PLAYER_MOVE"
	Payload     interface{}              // Type-specific payload data
	Timestamp   int64
}
