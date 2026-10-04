package game

import "errors"

type Board struct {
	cells [3][3]rune
}

// NewBoard creates an empty tic-tac-toe board.
func NewBoard() *Board {
	return &Board{}
}

// PlaceMove puts a symbol in an empty, in-bounds cell.
func (b *Board) PlaceMove(row int, col int, symbol rune) error {
	if row < 0 || row >= len(b.cells) || col < 0 || col >= len(b.cells[row]) {
		return errors.New("position is out of bounds")
	}
	if b.cells[row][col] != 0 {
		return errors.New("position is already occupied")
	}

	b.cells[row][col] = symbol
	return nil
}

// CellAt returns the symbol at an in-bounds position.
func (b *Board) CellAt(row int, col int) (rune, error) {
	if row < 0 || row >= len(b.cells) || col < 0 || col >= len(b.cells[row]) {
		return 0, errors.New("position is out of bounds")
	}

	return b.cells[row][col], nil
}

// IsFull reports whether every cell is occupied.
func (b *Board) IsFull() bool {
	for _, row := range b.cells {
		for _, cell := range row {
			if cell == 0 {
				return false
			}
		}
	}

	return true
}

// CheckWinner returns the symbol that has a complete line, if any.
func (b *Board) CheckWinner() (rune, bool) {
	lines := [8][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}

	for _, line := range lines {
		// get first value
		symbol := b.cells[line[0][0]][line[0][1]]
		if symbol == 0 {
			continue
		}
		// compare with other two of that same row /column of an lines
		if symbol == b.cells[line[1][0]][line[1][1]] &&
			symbol == b.cells[line[2][0]][line[2][1]] {
			return symbol, true
		}
	}

	return 0, false
}
