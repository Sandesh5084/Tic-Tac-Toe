package game

import "errors"

type GameState int

const (
	Playing GameState = iota
	Won
	Draw
)

type Game struct {
	board         *Board
	players       [2]*Player
	currentPlayer int
	state         GameState
}

// NewGame creates a match with an empty board and Player 1 to move.
func NewGame(player1 *Player, player2 *Player) *Game {
	return &Game{
		board:         NewBoard(),
		players:       [2]*Player{player1, player2},
		currentPlayer: 0,
		state:         Playing,
	}
}

// CellAt returns the value at the requested board position.
func (g *Game) CellAt(row int, col int) (rune, error) {
	return g.board.CellAt(row, col)
}

// CurrentPlayer returns a copy of the player whose turn it is.
func (g *Game) CurrentPlayer() *Player {
	player := g.players[g.currentPlayer]
	if player == nil {
		return nil
	}

	playerCopy := *player
	return &playerCopy
}

// State returns whether the match is playing, won, or drawn.
func (g *Game) State() GameState {
	return g.state
}

// Restart resets the board, state, and turn while keeping the same players.
func (g *Game) Restart() {
	g.board = NewBoard()
	g.currentPlayer = 0
	g.state = Playing
}

// MakeMove places the current player's symbol and updates the game status.
func (g *Game) MakeMove(move Move) error {
	if g.state != Playing {
		return errors.New("game is already finished")
	}

	player := g.players[g.currentPlayer]
	if err := g.board.PlaceMove(move.Row, move.Column, player.Symbol); err != nil {
		return err
	}

	if _, won := g.board.CheckWinner(); won {
		g.state = Won
		return nil
	}
	if g.board.IsFull() {
		g.state = Draw
		return nil
	}

	g.currentPlayer = (g.currentPlayer + 1) % len(g.players)
	return nil
}
