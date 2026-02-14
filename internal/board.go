package internal

import (
	"fmt"
	"strings"
)

type Symbol string

const (
	SymbolX     Symbol = "X"
	SymbolO     Symbol = "O"
	SymbolEmpty Symbol = ""
)

type BoardCell struct {
	Symbol Symbol
}

type Board struct {
	cells [9]BoardCell
}

func NewBoard() *Board {
	board := &Board{}
	for i := 0; i < 9; i++ {
		board.cells[i] = BoardCell{Symbol: SymbolEmpty}
	}
	return board
}

func (b *Board) SetCell(index int, symbol Symbol) {
	if index >= 0 && index < 9 {
		b.cells[index].Symbol = symbol
	}
}

func (b *Board) GetCell(index int) Symbol {
	if index >= 0 && index < 9 {
		return b.cells[index].Symbol
	}
	return SymbolEmpty
}

func (b *Board) GetCells() [9]BoardCell {
	return b.cells
}

func (b *Board) Clear() {
	for i := 0; i < 9; i++ {
		b.cells[i].Symbol = SymbolEmpty
	}
}

func (b *Board) String() string {
	var sb strings.Builder
	sb.WriteString("\n")

	for i := 0; i < 9; i++ {
		symbol := b.cells[i].Symbol
		if symbol == SymbolEmpty {
			sb.WriteString(" . ")
		} else {
			sb.WriteString(fmt.Sprintf(" %s ", symbol))
		}

		if i%3 == 2 { // End of row
			sb.WriteString("\n")
			if i < 6 { // Not the last row
				sb.WriteString("---|---|---\n")
			}
		} else {
			sb.WriteString("|")
		}
	}

	return sb.String()
}

func (b *Board) ToArray() []string {
	result := make([]string, 9)
	for i := 0; i < 9; i++ {
		result[i] = string(b.cells[i].Symbol)
	}
	return result
}

func (b *Board) FromArray(symbols []string) {
	for i := 0; i < 9 && i < len(symbols); i++ {
		b.cells[i].Symbol = Symbol(symbols[i])
	}
}
