package internal

import "fmt"

type GameStatus string

const (
	StatusInProgress GameStatus = "IN_PROGRESS"
	StatusWon        GameStatus = "WON"
	StatusDraw       GameStatus = "DRAW"
)

type Game struct {
	board       *Board
	currentTurn Symbol
	status      GameStatus
	winner      Symbol
}

func NewGame() *Game {
	return &Game{
		board:       NewBoard(),
		currentTurn: SymbolX,
		status:      StatusInProgress,
		winner:      SymbolEmpty,
	}
}

func (g *Game) MakeMove(row, col int, symbol Symbol) error {
	//if !ValidatePosition(row, col) {
	//	return fmt.Errorf("invalid position: row=%d, col=%d", row, col)
	//}

	if symbol != g.currentTurn {
		return fmt.Errorf("not your turn, current turn: %s", g.currentTurn)
	}

	if g.status != StatusInProgress {
		return fmt.Errorf("game is already over")
	}
	//
	//index := PositionToIndex(row, col)
	//
	//if g.board.GetCell(index) != SymbolEmpty {
	//	return fmt.Errorf("cell already occupied")
	//}
	//
	//g.board.SetCell(index, symbol)
	//g.checkGameState()
	//
	//if g.status == StatusInProgress {
	//	g.switchTurn()
	//}

	return nil
}

func (g *Game) GetBoard() *Board { return g.board }
