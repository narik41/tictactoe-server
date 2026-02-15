package internal

import (
	"fmt"
	"sync"
	"time"
)

type GameSessionStatus string

const (
	SessionWaitingForPlayers GameSessionStatus = "WAITING_FOR_PLAYERS"
	SessionReady             GameSessionStatus = "READY"
	SessionInProgress        GameSessionStatus = "IN_PROGRESS"
	SessionCompleted         GameSessionStatus = "COMPLETED"
	SessionAbandoned         GameSessionStatus = "ABANDONED"
)

type GameSession struct {
	Id        string
	Game      *Game
	PlayerX   *PlayerInfo
	PlayerO   *PlayerInfo
	Status    GameSessionStatus
	CreatedAt time.Time
	StartedAt time.Time
	EndedAt   time.Time
	mu        sync.RWMutex
}

type PlayerInfo struct {
	SessionID string // Reference to the connection session
	Username  string
	Symbol    Symbol // X or O
	IsReady   bool
	//client    *Client
	MyTurn bool
}

func NewGameSession(sessionID string) *GameSession {
	return &GameSession{
		Id:        sessionID,
		Game:      NewGame(),
		Status:    SessionWaitingForPlayers,
		CreatedAt: time.Now(),
	}
}

func (gs *GameSession) AddPlayer(sessionID, username string) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.PlayerX != nil && gs.PlayerO != nil {
		return fmt.Errorf("game session is full")
	}

	if gs.PlayerX == nil {
		gs.PlayerX = &PlayerInfo{
			SessionID: sessionID,
			Username:  username,
			Symbol:    SymbolX,
			IsReady:   false,
			//client:    client,
			MyTurn: false,
		}
	} else if gs.PlayerO == nil {
		gs.PlayerO = &PlayerInfo{
			SessionID: sessionID,
			Username:  username,
			Symbol:    SymbolO,
			IsReady:   false,
			MyTurn:    true,
			//client:    client,
		}
	}

	if gs.PlayerX != nil && gs.PlayerO != nil {
		gs.Status = SessionReady
	}

	return nil
}

func (gs *GameSession) RemovePlayer(sessionID string) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.PlayerX != nil && gs.PlayerX.SessionID == sessionID {
		gs.PlayerX = nil
		gs.handlePlayerDisconnect()
		return nil
	}

	if gs.PlayerO != nil && gs.PlayerO.SessionID == sessionID {
		gs.PlayerO = nil
		gs.handlePlayerDisconnect()
		return nil
	}

	return fmt.Errorf("player not found in session")
}

func (gs *GameSession) handlePlayerDisconnect() {
	if gs.Status == SessionInProgress {
		gs.Status = SessionAbandoned
		gs.EndedAt = time.Now()
	} else {
		gs.Status = SessionWaitingForPlayers
	}
}

func (gs *GameSession) Start() error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Status != SessionReady {
		return fmt.Errorf("cannot start game, status: %s", gs.Status)
	}

	if gs.PlayerX == nil || gs.PlayerO == nil {
		return fmt.Errorf("both players must be present")
	}

	gs.Status = SessionInProgress
	gs.StartedAt = time.Now()
	// relay the message

	return nil
}

func (gs *GameSession) MakeMove(sessionID string, position int) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Status != SessionInProgress {
		return fmt.Errorf("game is not in progress")
	}

	playerSymbol, err := gs.getPlayerSymbol(sessionID)
	if err != nil {
		return err
	}

	if err := gs.Game.MakeMove(position, playerSymbol); err != nil {
		return err
	}

	if gs.Game.IsGameOver() {
		gs.Status = SessionCompleted
		gs.EndedAt = time.Now()
	}

	return nil
}

func (g *Game) IsGameOver() bool {
	return g.status == StatusWon || g.status == StatusDraw
}

func (gs *GameSession) getPlayerSymbol(sessionID string) (Symbol, error) {
	if gs.PlayerX != nil && gs.PlayerX.SessionID == sessionID {
		return SymbolX, nil
	}
	if gs.PlayerO != nil && gs.PlayerO.SessionID == sessionID {
		return SymbolO, nil
	}
	return SymbolEmpty, fmt.Errorf("player not in this game session")
}

func (gs *GameSession) GetPlayerInfo(sessionID string) (*PlayerInfo, error) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	if gs.PlayerX != nil && gs.PlayerX.SessionID == sessionID {
		return gs.PlayerX, nil
	}
	if gs.PlayerO != nil && gs.PlayerO.SessionID == sessionID {
		return gs.PlayerO, nil
	}
	return nil, fmt.Errorf("player not found")
}

func (gs *GameSession) GetOpponentInfo(sessionID string) (*PlayerInfo, error) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	if gs.PlayerX != nil && gs.PlayerX.SessionID == sessionID {
		if gs.PlayerO != nil {
			return gs.PlayerO, nil
		}
		return nil, fmt.Errorf("opponent not found")
	}
	if gs.PlayerO != nil && gs.PlayerO.SessionID == sessionID {
		if gs.PlayerX != nil {
			return gs.PlayerX, nil
		}
		return nil, fmt.Errorf("opponent not found")
	}
	return nil, fmt.Errorf("player not in this game")
}

// IsPlayerTurn checks if it's the given player's turn
func (gs *GameSession) IsPlayerTurn(sessionID string) bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	_, err := gs.getPlayerSymbol(sessionID)
	if err != nil {
		return false
	}
	return false
	//return gs.Game.GetCurrentTurn() == playerSymbol
}

// GetBothPlayerSessionIDs returns both player session IDs
func (gs *GameSession) GetBothPlayerSessionIDs() []string {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	var sessionIDs []string
	if gs.PlayerX != nil {
		sessionIDs = append(sessionIDs, gs.PlayerX.SessionID)
	}
	if gs.PlayerO != nil {
		sessionIDs = append(sessionIDs, gs.PlayerO.SessionID)
	}
	return sessionIDs
}

func (gs *GameSession) IsFull() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.PlayerX != nil && gs.PlayerO != nil
}

func (gs *GameSession) HasPlayer(sessionID string) bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	if gs.PlayerX != nil && gs.PlayerX.SessionID == sessionID {
		return true
	}
	if gs.PlayerO != nil && gs.PlayerO.SessionID == sessionID {
		return true
	}
	return false
}
