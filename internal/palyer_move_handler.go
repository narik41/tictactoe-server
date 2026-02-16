package internal

import (
	"encoding/json"
	"log"

	"github.com/narik41/tictactoe-helper/core"
	"github.com/narik41/tictactoe-server/internal/decoder"
)

type PlayerMoveHandler struct {
	gameSessionManager *GameSessionManager
}

func NewPlayerMoveHandler(gameSessionManager *GameSessionManager) PlayerMoveHandler {
	return PlayerMoveHandler{
		gameSessionManager: gameSessionManager,
	}
}

func (a PlayerMoveHandler) Handle(msg *decoder.DecodedMessage, sessionId string) (*HandlerResponse, error) {
	log.Println("PlayerMoveHandler.Handle")
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var loginPayload core.Version1PositionMoveRequestPayload
	if err := json.Unmarshal(jsonBytes, &loginPayload); err != nil {
		return nil, err
	}

	gameSession, err := a.gameSessionManager.GetSessionByPlayer(sessionId)
	if err != nil {
		return nil, err
	}

	err = gameSession.MakeMove(sessionId, loginPayload.Position)
	if err != nil {
		return nil, err
	}

	if gameSession.Game.IsGameEnd() {
		return &HandlerResponse{
			Broadcast: true,
			Recipients: []string{
				gameSession.PlayerO.SessionID, gameSession.PlayerX.SessionID,
			},
			MessageType: core.GAME_END,
			Payload: &core.Version1GameEndPayload{
				Result: "",
				Winner: string(gameSession.Game.GetWinnerSymbol()),
			},
		}, nil
	}

	return &HandlerResponse{
		Broadcast: true,
		Recipients: []string{
			gameSession.PlayerO.SessionID, gameSession.PlayerX.SessionID,
		},
		MessageType: core.PLAYER_MOVE_RESPONSE,
		Payload: &core.Version1PositionMovedResponsePayload{
			MovedByUser:     loginPayload.Symbol,
			MovedToPosition: loginPayload.Position,
			TurnSymbol:      string(gameSession.Game.GetCurrentTurn()),
		},
	}, nil
}

func (a PlayerMoveHandler) RequiredStates() []SessionState {
	return []SessionState{
		LoggedIn,
		IN_GAME,
	}
}
