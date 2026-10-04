package game

type Player struct {
	Name   string
	Symbol rune
}

// NewPlayer creates a player with the given name and symbol.
func NewPlayer(name string, symbol rune) *Player {
	return &Player{
		Name:   name,
		Symbol: symbol,
	}
}
