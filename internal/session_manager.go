package internal

import (
	"fmt"
	"sync"
)

type GameSessionManager struct {
	sessions        map[string]*GameSession
	playerToSession map[string]string
}

type SessionManager struct {
	sessions   map[string]*Session // sessionId -> Session
	byUsername map[string]*Session // username -> Session
	mu         sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:   make(map[string]*Session),
		byUsername: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession(client *Client) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		Id:        UUID("sess"),
		Client:    client,
		State:     Guest,
		CreatedAt: GetNPTToUtcInMillisecond(),
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

func (sm *SessionManager) GetSessionByUsername(username string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.byUsername[username]
	return session, exists
}

func (sm *SessionManager) RegisterUsername(sessionID, username string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	// Check if username already taken
	if _, taken := sm.byUsername[username]; taken {
		return fmt.Errorf("username already in use")
	}

	session.Username = username
	sm.byUsername[username] = session

	return nil
}

func (sm *SessionManager) RemoveSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return
	}

	if session.Username != "" {
		delete(sm.byUsername, session.Username)
	}

	delete(sm.sessions, sessionID)
}

func (sm *SessionManager) GetAllSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}
