package game

import "testing"

// TestNewGame checks the initial state of a newly created game.
func TestNewGame(t *testing.T) {
	player1 := NewPlayer("Sandesh", 'X')
	player2 := NewPlayer("Rahul", 'O')

	game := NewGame(player1, player2)

	if game.board == nil {
		t.Fatal("NewGame() did not initialize the board")
	}
	if game.board.IsFull() {
		t.Error("NewGame() initialized a non-empty board")
	}
	if game.players[0] != player1 {
		t.Errorf("first player = %p; want %p", game.players[0], player1)
	}
	if game.players[1] != player2 {
		t.Errorf("second player = %p; want %p", game.players[1], player2)
	}
	if game.currentPlayer != 0 {
		t.Errorf("currentPlayer = %d; want 0", game.currentPlayer)
	}
	if game.state != Playing {
		t.Errorf("state = %v; want Playing", game.state)
	}
}

// TestRestart starts a fresh match with the same players.
func TestRestart(t *testing.T) {
	player1 := NewPlayer("Sandesh", 'X')
	player2 := NewPlayer("Rahul", 'O')
	game := NewGame(player1, player2)
	if err := game.MakeMove(Move{Row: 1, Column: 1}); err != nil {
		t.Fatalf("MakeMove() returned an error: %v", err)
	}

	game.Restart()

	if game.State() != Playing {
		t.Errorf("State() after Restart() = %v; want Playing", game.State())
	}
	if game.CurrentPlayer().Symbol != player1.Symbol {
		t.Errorf("CurrentPlayer() after Restart() = %q; want %q", game.CurrentPlayer().Symbol, player1.Symbol)
	}
	if cell, err := game.CellAt(1, 1); err != nil {
		t.Fatalf("CellAt() returned an error: %v", err)
	} else if cell != 0 {
		t.Errorf("cell after Restart() = %q; want empty", cell)
	}
}

// TestGameReadOnlyAccessors checks cell, player, and state getters.
func TestGameReadOnlyAccessors(t *testing.T) {
	player1 := NewPlayer("Sandesh", 'X')
	player2 := NewPlayer("Rahul", 'O')
	game := NewGame(player1, player2)

	if got, err := game.CellAt(0, 0); err != nil {
		t.Fatalf("CellAt() returned an error: %v", err)
	} else if got != 0 {
		t.Errorf("CellAt(0, 0) = %q; want empty", got)
	}
	if _, err := game.CellAt(-1, 0); err == nil {
		t.Error("CellAt() returned nil error for an invalid position")
	}
	if got := game.CurrentPlayer(); got == nil || got.Name != player1.Name || got.Symbol != player1.Symbol {
		t.Errorf("CurrentPlayer() = %+v; want a copy of %+v", got, player1)
	} else {
		got.Symbol = 'O'
		if player1.Symbol != 'X' {
			t.Errorf("changing CurrentPlayer() result changed stored player symbol to %q", player1.Symbol)
		}
	}
	if got := game.State(); got != Playing {
		t.Errorf("State() = %v; want Playing", got)
	}

	if err := game.MakeMove(Move{Row: 0, Column: 0}); err != nil {
		t.Fatalf("MakeMove() returned an error: %v", err)
	}
	if got := game.CurrentPlayer(); got == nil || got.Name != player2.Name || got.Symbol != player2.Symbol {
		t.Errorf("CurrentPlayer() after X moves = %+v; want %+v", got, player2)
	}
}

// TestMakeMove checks that a valid move is placed and advances the turn.
func TestMakeMove(t *testing.T) {
	player1 := NewPlayer("Sandesh", 'X')
	player2 := NewPlayer("Rahul", 'O')
	game := NewGame(player1, player2)

	if err := game.MakeMove(Move{Row: 0, Column: 0}); err != nil {
		t.Fatalf("MakeMove() returned an error: %v", err)
	}
	if got := game.board.cells[0][0]; got != player1.Symbol {
		t.Errorf("board cell = %q; want player symbol %q", got, player1.Symbol)
	}
	if game.currentPlayer != 1 {
		t.Errorf("currentPlayer = %d after valid move; want 1", game.currentPlayer)
	}
}

// TestMakeMoveInvalidPositionDoesNotSwitchTurn checks invalid moves preserve the turn.
func TestMakeMoveInvalidPositionDoesNotSwitchTurn(t *testing.T) {
	game := NewGame(NewPlayer("Sandesh", 'X'), NewPlayer("Rahul", 'O'))

	if err := game.MakeMove(Move{Row: 3, Column: 0}); err == nil {
		t.Fatal("MakeMove() returned nil for an invalid position; want an error")
	}
	if game.currentPlayer != 0 {
		t.Errorf("currentPlayer = %d after invalid move; want 0", game.currentPlayer)
	}
}

// TestMakeMoveOccupiedPositionDoesNotSwitchTurn checks occupied moves preserve the turn.
func TestMakeMoveOccupiedPositionDoesNotSwitchTurn(t *testing.T) {
	game := NewGame(NewPlayer("Sandesh", 'X'), NewPlayer("Rahul", 'O'))

	if err := game.MakeMove(Move{Row: 1, Column: 1}); err != nil {
		t.Fatalf("first MakeMove() returned an error: %v", err)
	}
	if err := game.MakeMove(Move{Row: 1, Column: 1}); err == nil {
		t.Fatal("MakeMove() returned nil for an occupied position; want an error")
	}
	if game.currentPlayer != 1 {
		t.Errorf("currentPlayer = %d after occupied move; want 1", game.currentPlayer)
	}
	if got := game.board.cells[1][1]; got != 'X' {
		t.Errorf("occupied cell = %q; want %q", got, 'X')
	}
}

// TestMakeMoveWinningMoveSetsWonAndDoesNotSwitchTurn checks the winning move behavior.
func TestMakeMoveWinningMoveSetsWonAndDoesNotSwitchTurn(t *testing.T) {
	game := NewGame(NewPlayer("Sandesh", 'X'), NewPlayer("Rahul", 'O'))
	moves := []Move{
		{Row: 0, Column: 0},
		{Row: 1, Column: 0},
		{Row: 0, Column: 1},
		{Row: 1, Column: 1},
		{Row: 0, Column: 2},
	}

	for _, move := range moves {
		if err := game.MakeMove(move); err != nil {
			t.Fatalf("MakeMove(%+v) returned an error: %v", move, err)
		}
	}

	if game.State() != Won {
		t.Errorf("State() = %v after winning move; want Won", game.State())
	}
	if game.currentPlayer != 0 {
		t.Errorf("currentPlayer = %d after winning move; want 0", game.currentPlayer)
	}
}

// TestMakeMoveDrawDoesNotSwitchTurn checks that a full, tied board becomes a draw.
func TestMakeMoveDrawDoesNotSwitchTurn(t *testing.T) {
	game := NewGame(NewPlayer("Sandesh", 'X'), NewPlayer("Rahul", 'O'))
	moves := []Move{
		{Row: 0, Column: 0},
		{Row: 1, Column: 1},
		{Row: 2, Column: 2},
		{Row: 0, Column: 1},
		{Row: 0, Column: 2},
		{Row: 2, Column: 0},
		{Row: 1, Column: 0},
		{Row: 1, Column: 2},
		{Row: 2, Column: 1},
	}

	for _, move := range moves {
		if err := game.MakeMove(move); err != nil {
			t.Fatalf("MakeMove(%+v) returned an error: %v", move, err)
		}
	}

	if game.State() != Draw {
		t.Errorf("State() = %v after draw; want Draw", game.State())
	}
	if game.currentPlayer != 0 {
		t.Errorf("currentPlayer = %d after draw; want 0", game.currentPlayer)
	}
}

// TestMakeMoveRejectsMovesAfterWin checks that a won game rejects further moves.
func TestMakeMoveRejectsMovesAfterWin(t *testing.T) {
	game := NewGame(NewPlayer("Sandesh", 'X'), NewPlayer("Rahul", 'O'))
	for _, move := range []Move{
		{Row: 0, Column: 0},
		{Row: 1, Column: 0},
		{Row: 0, Column: 1},
		{Row: 1, Column: 1},
		{Row: 0, Column: 2},
	} {
		if err := game.MakeMove(move); err != nil {
			t.Fatalf("MakeMove(%+v) returned an error: %v", move, err)
		}
	}

	if err := game.MakeMove(Move{Row: 2, Column: 2}); err == nil {
		t.Fatal("MakeMove() returned nil after a win; want an error")
	}
	if game.currentPlayer != 0 {
		t.Errorf("currentPlayer = %d after rejected move; want 0", game.currentPlayer)
	}
	if got := game.board.cells[2][2]; got != 0 {
		t.Errorf("cell after rejected move = %q; want empty", got)
	}
}

// TestMakeMoveRejectsMovesAfterDraw checks that a drawn game rejects further moves.
func TestMakeMoveRejectsMovesAfterDraw(t *testing.T) {
	game := NewGame(NewPlayer("Sandesh", 'X'), NewPlayer("Rahul", 'O'))
	for _, move := range []Move{
		{Row: 0, Column: 0},
		{Row: 1, Column: 1},
		{Row: 2, Column: 2},
		{Row: 0, Column: 1},
		{Row: 0, Column: 2},
		{Row: 2, Column: 0},
		{Row: 1, Column: 0},
		{Row: 1, Column: 2},
		{Row: 2, Column: 1},
	} {
		if err := game.MakeMove(move); err != nil {
			t.Fatalf("MakeMove(%+v) returned an error: %v", move, err)
		}
	}

	if err := game.MakeMove(Move{Row: 2, Column: 2}); err == nil {
		t.Fatal("MakeMove() returned nil after a draw; want an error")
	} else if err.Error() != "game is already finished" {
		t.Errorf("MakeMove() error = %q after a draw; want game-finished error", err)
	}
	if game.currentPlayer != 0 {
		t.Errorf("currentPlayer = %d after rejected move; want 0", game.currentPlayer)
	}
	if got := game.board.cells[2][2]; got != 'X' {
		t.Errorf("cell after rejected move = %q; want %q", got, 'X')
	}
}

// TestMakeMoveKeepsPlayingAndSwitchesTurn checks non-terminal turn progression.
func TestMakeMoveKeepsPlayingAndSwitchesTurn(t *testing.T) {
	game := NewGame(NewPlayer("Sandesh", 'X'), NewPlayer("Rahul", 'O'))

	if err := game.MakeMove(Move{Row: 0, Column: 0}); err != nil {
		t.Fatalf("MakeMove() returned an error: %v", err)
	}

	if game.state != Playing {
		t.Errorf("state = %v after a non-winning move; want Playing", game.state)
	}
	if game.currentPlayer != 1 {
		t.Errorf("currentPlayer = %d after a non-winning move; want 1", game.currentPlayer)
	}
}
