package main

import (
	"fmt"

	"myapp/internal/cli"
	"myapp/internal/game"
)

// main starts the application.
func main() {
	fmt.Println("Welcome to Tic-Tac-Toe!")

	player1 := game.NewPlayer("Player 1", 'X')
	player2 := game.NewPlayer("Player 2", 'O')
	match := game.NewGame(player1, player2)

	cli.NewCLI(match).Run()
}
