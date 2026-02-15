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

func (g *Game) MakeMove(position int, symbol Symbol) error {

	if position < 0 || position > 8 {
		return fmt.Errorf("invalid position: position ", position)
	}

	if symbol != g.currentTurn {
		return fmt.Errorf("not your turn, current turn: %s", g.currentTurn)
	}

	if g.status != StatusInProgress {
		return fmt.Errorf("game is already over")
	}

	if g.board.GetCell(position) != SymbolEmpty {
		return fmt.Errorf("cell already occupied")
	}

	g.board.SetCell(position, symbol)
	g.checkGameState()

	if g.status == StatusInProgress {
		g.switchTurn()
	}

	return nil
}

func (g *Game) checkGameState() {

	if winner := g.checkWinner(); winner != SymbolEmpty {
		g.status = StatusWon
		g.winner = winner
		return
	}

	if g.isBoardFull() {
		g.status = StatusDraw
		return
	}
	g.status = StatusInProgress
}

func (g *Game) checkWinner() Symbol {

	winningCombos := [][]int{
		{0, 1, 2}, // Top row
		{3, 4, 5}, // Middle row
		{6, 7, 8}, // Bottom row
		{0, 3, 6}, // Left column
		{1, 4, 7}, // Middle column
		{2, 5, 8}, // Right column
		{0, 4, 8}, // Diagonal top-left to bottom-right
		{2, 4, 6}, // Diagonal top-right to bottom-left
	}

	for _, combo := range winningCombos {
		first := g.board.GetCell(combo[0])
		if first != SymbolEmpty &&
			first == g.board.GetCell(combo[1]) &&
			first == g.board.GetCell(combo[2]) {
			return first
		}
	}

	return SymbolEmpty
}

func (g *Game) isBoardFull() bool {
	for i := 0; i < 9; i++ {
		if g.board.GetCell(i) == SymbolEmpty {
			return false
		}
	}
	return true
}

func (g *Game) switchTurn() {
	if g.currentTurn == SymbolX {
		g.currentTurn = SymbolO
	} else {
		g.currentTurn = SymbolX
	}
}

func (g *Game) GetBoard() *Board { return g.board }

func (g *Game) GetCurrentTurn() Symbol { return g.currentTurn }

func (g *Game) IsGameEnd() bool         { return g.status == StatusWon }
func (g *Game) GetWinnerSymbol() Symbol { return g.winner }
