package game

import "testing"

// TestMoveStoresPosition checks that a move stores its row and column.
func TestMoveStoresPosition(t *testing.T) {
	move := Move{Row: 1, Column: 2}

	if move.Row != 1 {
		t.Errorf("Row = %d; want 1", move.Row)
	}
	if move.Column != 2 {
		t.Errorf("Column = %d; want 2", move.Column)
	}
}
