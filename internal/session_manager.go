package internal

import (
	"sync"

	"github.com/narik41/tictactoe-helper/core"
)

type SessionManager struct {
	sessions map[string]*Session // sessionId -> Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession(client *Client) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		Id:        core.UUID("player"),
		Client:    client,
		State:     Guest,
		CreatedAt: core.GetNPTToUtcInMillisecond(),
	}

	sm.sessions[session.Id] = session
	return session
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	return session, exists
}
