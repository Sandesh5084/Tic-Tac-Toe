package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"myapp/internal/game"
)

var ErrQuit = errors.New("quit requested")

type CLI struct {
	game   *game.Game
	input  *bufio.Reader
	output io.Writer
}

// NewCLI creates a CLI connected to the game and standard input/output.
func NewCLI(game *game.Game) *CLI {
	return &CLI{
		game:   game,
		input:  bufio.NewReader(os.Stdin),
		output: os.Stdout,
	}
}

// ReadMove prompts the current player and reads a row and column.
func (c *CLI) ReadMove() (game.Move, error) {
	player := c.game.CurrentPlayer()
	if player == nil {
		return game.Move{}, fmt.Errorf("cannot read move: current player is nil")
	}
	if _, err := fmt.Fprintf(c.output, "Player %c's turn\nEnter row and column: ", player.Symbol); err != nil {
		return game.Move{}, fmt.Errorf("write move prompt: %w", err)
	}

	line, err := c.input.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		return game.Move{}, fmt.Errorf("read row and column: %w", err)
	}

	if strings.TrimSpace(line) == "q" {
		return game.Move{}, ErrQuit
	}

	fields := strings.Fields(line)
	if len(fields) != 2 {
		return game.Move{}, fmt.Errorf("expected two integers")
	}

	row, err := strconv.Atoi(fields[0])
	if err != nil {
		return game.Move{}, fmt.Errorf("expected two integers: invalid row %q", fields[0])
	}
	column, err := strconv.Atoi(fields[1])
	if err != nil {
		return game.Move{}, fmt.Errorf("expected two integers: invalid column %q", fields[1])
	}

	return game.Move{Row: row, Column: column}, nil
}

// Run coordinates input, moves, board display, and the final result.
func (c *CLI) Run() {
	fmt.Fprintln(c.output, "Enter q + Enter at any time to quit.")

	for {
		c.PrintBoard()

		for c.game.State() == game.Playing {
			move, err := c.ReadMove()
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, ErrQuit) {
					return
				}
				fmt.Fprintf(c.output, "Invalid input: %v\n", err)
				continue
			}

			if err := c.game.MakeMove(move); err != nil {
				fmt.Fprintf(c.output, "Invalid move: %v\n", err)
				continue
			}

			c.PrintBoard()
		}

		switch c.game.State() {
		case game.Won:
			player := c.game.CurrentPlayer()
			if player != nil {
				fmt.Fprintf(c.output, "Player %s (%c) wins!\n", player.Name, player.Symbol)
			}
		case game.Draw:
			fmt.Fprintln(c.output, "It's a draw!")
		}

		retry, err := c.readReplayChoice()
		if err != nil || retry == 'q' {
			return
		}
		c.game.Restart()
	}
}

// readReplayChoice prompts for another match or an exit.
func (c *CLI) readReplayChoice() (rune, error) {
	for {
		if _, err := fmt.Fprint(c.output, "Press r + Enter to retry or q + Enter to quit: "); err != nil {
			return 0, fmt.Errorf("write replay prompt: %w", err)
		}

		line, err := c.input.ReadString('\n')
		if err != nil && !(errors.Is(err, io.EOF) && line != "") {
			return 0, fmt.Errorf("read replay choice: %w", err)
		}

		switch strings.ToLower(strings.TrimSpace(line)) {
		case "r":
			return 'r', nil
		case "q":
			return 'q', nil
		default:
			if _, writeErr := fmt.Fprintln(c.output, "Invalid choice. Enter r to retry or q to quit."); writeErr != nil {
				return 0, fmt.Errorf("write invalid-choice message: %w", writeErr)
			}
		}
	}
}

// PrintBoard displays the current board as a three-by-three grid.
func (c *CLI) PrintBoard() {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			cell, err := c.game.CellAt(row, col)
			if err != nil {
				fmt.Fprintf(c.output, "Unable to read board cell (%d, %d): %v\n", row, col, err)
				return
			}
			if cell == 0 {
				cell = ' '
			}

			fmt.Fprintf(c.output, " %c ", cell)
			if col < 2 {
				fmt.Fprint(c.output, "|")
			}
		}
		fmt.Fprintln(c.output)

		if row < 2 {
			fmt.Fprintln(c.output, "---+---+---")
		}
	}
}
