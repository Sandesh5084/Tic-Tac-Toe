package game

import "testing"

// TestNewBoard checks that a new board has only empty cells.
func TestNewBoard(t *testing.T) {
	board := NewBoard()

	if board == nil {
		t.Fatal("NewBoard() returned nil")
	}

	for row := range board.cells {
		for column, cell := range board.cells[row] {
			if cell != 0 {
				t.Errorf("cell at row %d, column %d = %q; want empty", row, column, cell)
			}
		}
	}
}

// TestPlaceMove checks that valid moves place their symbols.
func TestPlaceMove(t *testing.T) {
	tests := []struct {
		name   string
		row    int
		col    int
		symbol rune
	}{
		{name: "top-left", row: 0, col: 0, symbol: 'X'},
		{name: "bottom-right", row: 2, col: 2, symbol: 'O'},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			board := NewBoard()

			if err := board.PlaceMove(test.row, test.col, test.symbol); err != nil {
				t.Fatalf("PlaceMove() returned an error: %v", err)
			}
			if got := board.cells[test.row][test.col]; got != test.symbol {
				t.Errorf("cell at row %d, column %d = %q; want %q", test.row, test.col, got, test.symbol)
			}
		})
	}
}

// TestPlaceMovePositionAlreadyOccupied checks that occupied cells reject moves.
func TestPlaceMovePositionAlreadyOccupied(t *testing.T) {
	board := NewBoard()
	if err := board.PlaceMove(1, 1, 'X'); err != nil {
		t.Fatalf("first PlaceMove() returned an error: %v", err)
	}

	if err := board.PlaceMove(1, 1, 'O'); err == nil {
		t.Fatal("second PlaceMove() returned nil; want an error")
	}
	if got := board.cells[1][1]; got != 'X' {
		t.Errorf("occupied cell = %q; want %q", got, 'X')
	}
}

// TestPlaceMoveInvalidPosition checks that out-of-bounds moves are rejected.
func TestPlaceMoveInvalidPosition(t *testing.T) {
	tests := []struct {
		name string
		row  int
		col  int
	}{
		{name: "row below range", row: -1, col: 0},
		{name: "row above range", row: 3, col: 0},
		{name: "column below range", row: 0, col: -1},
		{name: "column above range", row: 0, col: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			board := NewBoard()
			if err := board.PlaceMove(test.row, test.col, 'X'); err == nil {
				t.Fatal("PlaceMove() returned nil; want an error")
			}
		})
	}
}

// TestIsFull_EmptyBoard checks that an empty board is not full.
func TestIsFull_EmptyBoard(t *testing.T) {
	board := NewBoard()

	if board.IsFull() {
		t.Error("IsFull() = true for an empty board; want false")
	}
}

// TestIsFull_PartiallyFilled checks that an incomplete board is not full.
func TestIsFull_PartiallyFilled(t *testing.T) {
	board := NewBoard()
	if err := board.PlaceMove(0, 0, 'X'); err != nil {
		t.Fatalf("PlaceMove() returned an error: %v", err)
	}

	if board.IsFull() {
		t.Error("IsFull() = true for a partially filled board; want false")
	}
}

// TestIsFull_FullBoard checks that a completely occupied board is full.
func TestIsFull_FullBoard(t *testing.T) {
	board := NewBoard()
	for row := range board.cells {
		for col := range board.cells[row] {
			symbol := 'X'
			if (row+col)%2 != 0 {
				symbol = 'O'
			}
			if err := board.PlaceMove(row, col, symbol); err != nil {
				t.Fatalf("PlaceMove(%d, %d, %q) returned an error: %v", row, col, symbol, err)
			}
		}
	}

	if !board.IsFull() {
		t.Error("IsFull() = false for a completely filled board; want true")
	}
}

// TestCheckWinner checks winning lines and a board with no winner.
func TestCheckWinner(t *testing.T) {
	type move struct {
		row, col int
		symbol   rune
	}

	tests := []struct {
		name   string
		moves  []move
		winner rune
		found  bool
	}{
		{
			name:   "horizontal",
			moves:  []move{{0, 0, 'X'}, {0, 1, 'X'}, {0, 2, 'X'}},
			winner: 'X',
			found:  true,
		},
		{
			name:   "vertical",
			moves:  []move{{0, 1, 'O'}, {1, 1, 'O'}, {2, 1, 'O'}},
			winner: 'O',
			found:  true,
		},
		{
			name:   "diagonal top-left to bottom-right",
			moves:  []move{{0, 0, 'X'}, {1, 1, 'X'}, {2, 2, 'X'}},
			winner: 'X',
			found:  true,
		},
		{
			name:   "diagonal top-right to bottom-left",
			moves:  []move{{0, 2, 'O'}, {1, 1, 'O'}, {2, 0, 'O'}},
			winner: 'O',
			found:  true,
		},
		{
			name:  "no winner",
			moves: []move{{0, 0, 'X'}, {0, 1, 'O'}, {1, 1, 'X'}, {2, 2, 'O'}},
			found: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			board := NewBoard()
			for _, move := range test.moves {
				if err := board.PlaceMove(move.row, move.col, move.symbol); err != nil {
					t.Fatalf("PlaceMove() returned an error: %v", err)
				}
			}

			winner, found := board.CheckWinner()
			if found != test.found {
				t.Fatalf("CheckWinner() found = %t; want %t", found, test.found)
			}
			if found && winner != test.winner {
				t.Errorf("CheckWinner() winner = %q; want %q", winner, test.winner)
			}
		})
	}
}
