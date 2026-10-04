package game

import "testing"

// TestNewPlayer checks that the constructor stores the player's details.
func TestNewPlayer(t *testing.T) {
	const name = "Sandesh"
	const symbol = 'X'

	player := NewPlayer(name, symbol)

	if player.Name != name {
		t.Errorf("Name = %q; want %q", player.Name, name)
	}
	if player.Symbol != symbol {
		t.Errorf("Symbol = %q; want %q", player.Symbol, symbol)
	}
}
