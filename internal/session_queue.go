package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/narik41/tictactoe-helper/core"
	"github.com/narik41/tictactoe-server/internal/game"
)

type SessionQueue struct {
	queue              []*Session
	gameSessionManager *GameSessionManager
	sender             *ResponseSender
	running            bool
}

func NewSessionQueue(gameSessionManager *GameSessionManager, sender *ResponseSender) *SessionQueue {
	mq := &SessionQueue{
		queue:              make([]*Session, 0),
		gameSessionManager: gameSessionManager,
		sender:             sender,
		running:            false,
	}
	return mq
}

func (mq *SessionQueue) Start() {
	if mq.running {
		//mq.mu.Unlock()
		return
	}
	mq.running = true

	go mq.matchmakingLoop()
}

func (mq *SessionQueue) Enqueue(session *Session) error {

	for _, s := range mq.queue {
		if s.Id == session.Id {
			return fmt.Errorf("session already in queue")
		}
	}

	mq.queue = append(mq.queue, session)
	session.State = WaitingForPair

	log.Printf("Session %s (%s) added to session queue. Queue size: %d",
		session.Id, session.Username, len(mq.queue))

	mq.sender.Send(session, &HandlerResponse{
		MessageType: core.WAITING_FOR_OPPONENT,
		Payload: map[string]interface{}{
			"message": "Waiting for an opponent...",
		},
	})

	return nil
}

func (mq *SessionQueue) Dequeue() *Session {

	if len(mq.queue) == 0 {
		return nil
	}

	session := mq.queue[0]
	mq.queue = mq.queue[1:]
	return session
}

func (mq *SessionQueue) Remove(sessionID string) error {

	for i, session := range mq.queue {
		if session.Id == sessionID {
			mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
			log.Printf("Session %s removed from matchmaking queue", sessionID)
			return nil
		}
	}

	return fmt.Errorf("session not in queue")
}

func (mq *SessionQueue) Size() int {
	return len(mq.queue)
}

func (mq *SessionQueue) matchmakingLoop() {
	log.Println("Session queue loop started")

	for {

		if !mq.running {

			break
		}

		queueSize := len(mq.queue)

		if queueSize >= 2 {
			mq.createMatch()
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}

	log.Println("Session queue loop stopped")
}

func (mq *SessionQueue) createMatch() {

	player1 := mq.Dequeue()
	player2 := mq.Dequeue()

	if player1 == nil || player2 == nil {
		if player1 != nil {
			mq.Enqueue(player1)
		}
		if player2 != nil {
			mq.Enqueue(player2)
		}
		return
	}

	log.Printf("Matching players: %s (%s) vs %s (%s)",
		player1.Id, player1.Username,
		player2.Id, player2.Username)

	// Create game session
	gameSession := mq.gameSessionManager.CreateSession()

	// Add both players to game
	err := mq.gameSessionManager.AddPlayerToSession(
		gameSession.Id,
		player1.Id,
		player1.Username,
	)
	if err != nil {
		log.Printf("Failed to add player1 to game: %v", err)
		mq.Enqueue(player1)
		mq.Enqueue(player2)
		return
	}

	err = mq.gameSessionManager.AddPlayerToSession(
		gameSession.Id,
		player2.Id,
		player2.Username,
	)
	if err != nil {
		log.Printf("Failed to add player2 to game: %v", err)
		mq.gameSessionManager.RemoveSession(gameSession.Id)
		mq.Enqueue(player1)
		mq.Enqueue(player2)
		return
	}

	player1.State = IN_GAME
	player2.State = IN_GAME
	gameSession.Start()

	log.Printf("Game %s started between %s and %s",
		gameSession.Game, player1.Username, player2.Username)

	mq.notifyGameStart(player1, player2, gameSession)
}

func (mq *SessionQueue) notifyGameStart(player1, player2 *Session, gameSession *game.GameSession) {

	playerXInfo, _ := gameSession.GetPlayerInfo(player1.Id)
	playerOInfo, _ := gameSession.GetPlayerInfo(player2.Id)

	mq.sender.Send(player1, &HandlerResponse{
		MessageType: core.GAME_START,
		Payload: &core.Version1GameStartPayload{
			GameId:     gameSession.Id,
			YourSymbol: string(playerXInfo.Symbol),
			YourTurn:   gameSession.Game.GetCurrentTurn() == playerXInfo.Symbol,
		},
	})

	// Notify Player 2
	mq.sender.Send(player2, &HandlerResponse{
		MessageType: core.GAME_START,
		Payload: &core.Version1GameStartPayload{
			GameId:     gameSession.Id,
			YourSymbol: string(playerOInfo.Symbol),
			YourTurn:   gameSession.Game.GetCurrentTurn() == playerOInfo.Symbol,
		},
	})

	log.Printf("Game start notifications sent to both players")
}
