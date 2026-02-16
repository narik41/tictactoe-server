package internal

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/narik41/tictactoe-helper/core"
	"github.com/narik41/tictactoe-server/internal/game"
)

type GameSessionManager struct {
	sessions        map[string]*game.GameSession // gameID -> GameSession
	playerToSession map[string]string            // playerSessionID -> gameID
	mu              sync.RWMutex
}

func NewGameSessionManager() *GameSessionManager {
	return &GameSessionManager{
		sessions:        make(map[string]*game.GameSession),
		playerToSession: make(map[string]string),
	}
}

func (gsm *GameSessionManager) CreateSession() *game.GameSession {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gameID := core.UUID("game")
	session := game.NewGameSession(gameID)
	gsm.sessions[gameID] = session

	log.Printf("Created game session: %s", gameID)
	return session
}

func (gsm *GameSessionManager) GetSession(gameID string) (*game.GameSession, error) {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	session, exists := gsm.sessions[gameID]
	if !exists {
		return nil, fmt.Errorf("game session not found")
	}
	return session, nil
}

func (gsm *GameSessionManager) GetSessionByPlayer(playerSessionID string) (*game.GameSession, error) {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	gameID, exists := gsm.playerToSession[playerSessionID]
	if !exists {
		return nil, fmt.Errorf("player not in any game session")
	}

	session, exists := gsm.sessions[gameID]
	if !exists {
		return nil, fmt.Errorf("game session not found")
	}

	return session, nil
}

func (gsm *GameSessionManager) AddPlayerToSession(gameID, playerSessionID, username string) error {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	session, exists := gsm.sessions[gameID]
	if !exists {
		return fmt.Errorf("game session not found")
	}

	if existingGameID, inGame := gsm.playerToSession[playerSessionID]; inGame {
		return fmt.Errorf("player already in game %s", existingGameID)
	}

	if err := session.AddPlayer(playerSessionID, username); err != nil {
		return err
	}

	gsm.playerToSession[playerSessionID] = gameID
	log.Printf("Added player %s (%s) to game %s", playerSessionID, username, gameID)
	return nil
}

func (gsm *GameSessionManager) RemovePlayerFromSession(playerSessionID string) error {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gameID, exists := gsm.playerToSession[playerSessionID]
	if !exists {
		return fmt.Errorf("player not in any game session")
	}

	session, exists := gsm.sessions[gameID]
	if !exists {
		return fmt.Errorf("game session not found")
	}

	// Remove player from session
	if err := session.RemovePlayer(playerSessionID); err != nil {
		return err
	}

	// Remove mapping
	delete(gsm.playerToSession, playerSessionID)

	log.Printf("Removed player %s from game %s", playerSessionID, gameID)

	// Clean up if game is abandoned or empty
	if session.Status == game.SessionAbandoned || !session.IsFull() {
		gsm.removeSession(gameID)
	}

	return nil
}

func (gsm *GameSessionManager) RemoveSession(gameID string) error {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	return gsm.removeSession(gameID)
}

func (gsm *GameSessionManager) removeSession(gameID string) error {
	session, exists := gsm.sessions[gameID]
	if !exists {
		return fmt.Errorf("game session not found")
	}

	// Remove player mappings
	if session.PlayerX != nil {
		delete(gsm.playerToSession, session.PlayerX.SessionID)
	}
	if session.PlayerO != nil {
		delete(gsm.playerToSession, session.PlayerO.SessionID)
	}

	// Remove session
	delete(gsm.sessions, gameID)

	log.Printf("Removed game session: %s", gameID)
	return nil
}

func (gsm *GameSessionManager) GetAllActiveSessions() []*game.GameSession {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	sessions := make([]*game.GameSession, 0, len(gsm.sessions))
	for _, session := range gsm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

func (gsm *GameSessionManager) GetSessionCount() int {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()
	return len(gsm.sessions)
}

func (gsm *GameSessionManager) IsPlayerInGame(playerSessionID string) bool {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	_, exists := gsm.playerToSession[playerSessionID]
	return exists
}

func (gsm *GameSessionManager) CleanupCompletedGames(maxAge time.Duration) int {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	now := time.Now()
	removed := 0

	for gameID, session := range gsm.sessions {
		shouldRemove := false

		if session.Status == game.SessionCompleted {
			if !session.EndedAt.IsZero() && now.Sub(session.EndedAt) > 5*time.Minute {
				shouldRemove = true
			}
		}

		if session.Status == game.SessionAbandoned {
			if !session.EndedAt.IsZero() && now.Sub(session.EndedAt) > 1*time.Minute {
				shouldRemove = true
			}
		}

		if now.Sub(session.CreatedAt) > maxAge {
			shouldRemove = true
		}

		if shouldRemove {
			gsm.removeSession(gameID)
			removed++
		}
	}

	if removed > 0 {
		log.Printf("Cleaned up %d old game sessions", removed)
	}

	return removed
}
