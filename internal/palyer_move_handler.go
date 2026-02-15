package internal

import (
	"encoding/json"
	"log"

	"github.com/narik41/tictactoe-message/core"
)

type PlayerMoveHandler struct {
	sessionManager *GameSessionManager
}

func NewPlayerMoveHandler(sessionManager *GameSessionManager) PlayerMoveHandler {
	return PlayerMoveHandler{
		sessionManager: sessionManager,
	}
}

func (a PlayerMoveHandler) Handle(msg *DecodedMessage, session *Session) (*HandlerResponse, error) {
	log.Println("PlayerMoveHandler.Handle")
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var loginPayload core.Version1PositionMoveRequestPayload
	if err := json.Unmarshal(jsonBytes, &loginPayload); err != nil {
		return nil, err
	}

	gameSession, err := a.sessionManager.GetSessionByPlayer(session.Id)
	if err != nil {
		return nil, err
	}

	err = gameSession.MakeMove(session.Id, loginPayload.Position)
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
