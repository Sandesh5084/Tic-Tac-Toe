package cli

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"myapp/internal/game"
)

// TestPrintBoard checks the rendered layout of the board.
func TestPrintBoard(t *testing.T) {
	match := game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O'))
	for _, move := range []game.Move{
		{Row: 0, Column: 0},
		{Row: 0, Column: 1},
		{Row: 1, Column: 1},
		{Row: 2, Column: 0},
		{Row: 2, Column: 2},
	} {
		if err := match.MakeMove(move); err != nil {
			t.Fatalf("MakeMove(%+v) returned an error: %v", move, err)
		}
	}

	cli := NewCLI(match)
	var output bytes.Buffer
	cli.output = &output
	cli.PrintBoard()

	const want = " X | O |   \n---+---+---\n   | X |   \n---+---+---\n O |   | X \n"
	if output.String() != want {
		t.Errorf("PrintBoard() output:\n%s\nwant:\n%s", output.String(), want)
	}
}

// TestReadMove checks that the CLI reads and prompts for a move.
func TestReadMove(t *testing.T) {
	match := game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O'))
	cli := NewCLI(match)
	cli.input = bufio.NewReader(strings.NewReader("1 2\n"))
	var output bytes.Buffer
	cli.output = &output

	move, err := cli.ReadMove()
	if err != nil {
		t.Fatalf("ReadMove() returned an error: %v", err)
	}
	if move != (game.Move{Row: 1, Column: 2}) {
		t.Errorf("ReadMove() = %+v; want {Row:1 Column:2}", move)
	}
	const wantPrompt = "Player X's turn\nEnter row and column: "
	if output.String() != wantPrompt {
		t.Errorf("prompt = %q; want %q", output.String(), wantPrompt)
	}
}

// TestReadMoveInvalidInput checks that non-integer input returns an error.
func TestReadMoveInvalidInput(t *testing.T) {
	cli := NewCLI(game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O')))
	cli.input = bufio.NewReader(strings.NewReader("1 two\n"))
	cli.output = io.Discard

	if _, err := cli.ReadMove(); err == nil {
		t.Fatal("ReadMove() returned nil error for non-integer input")
	}
}

// TestRunRetriesInvalidMoveAndReportsWinner checks move retries and win output.
func TestRunRetriesInvalidMoveAndReportsWinner(t *testing.T) {
	match := game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O'))
	cli := NewCLI(match)
	cli.input = bufio.NewReader(strings.NewReader("0 0\n0 0\n1 0\n0 1\n1 1\n0 2\n"))
	var output bytes.Buffer
	cli.output = &output

	cli.Run()

	if match.State() != game.Won {
		t.Fatalf("State() = %v after Run(); want Won", match.State())
	}
	if !strings.Contains(output.String(), "Invalid move: position is already occupied") {
		t.Errorf("Run() output does not report invalid move:\n%s", output.String())
	}
	if !strings.Contains(output.String(), "Player Sandesh (X) wins!") {
		t.Errorf("Run() output does not report the winner:\n%s", output.String())
	}
}

// TestRunReportsDraw checks that the CLI reports a drawn game.
func TestRunReportsDraw(t *testing.T) {
	match := game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O'))
	cli := NewCLI(match)
	cli.input = bufio.NewReader(strings.NewReader("0 0\n1 1\n2 2\n0 1\n0 2\n2 0\n1 0\n1 2\n2 1\n"))
	var output bytes.Buffer
	cli.output = &output

	cli.Run()

	if match.State() != game.Draw {
		t.Fatalf("State() = %v after Run(); want Draw", match.State())
	}
	if !strings.Contains(output.String(), "It's a draw!") {
		t.Errorf("Run() output does not report the draw:\n%s", output.String())
	}
}

// TestRunContinuesAfterMalformedLine checks malformed input is consumed before retrying.
func TestRunContinuesAfterMalformedLine(t *testing.T) {
	match := game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O'))
	cli := NewCLI(match)
	cli.input = bufio.NewReader(strings.NewReader("1 two\n0 0\n1 0\n0 1\n1 1\n0 2\n"))
	var output bytes.Buffer
	cli.output = &output

	cli.Run()

	if match.State() != game.Won {
		t.Fatalf("State() = %v after Run(); want Won", match.State())
	}
	if !strings.Contains(output.String(), "Invalid input: expected two integers: invalid column \"two\"") {
		t.Errorf("Run() output does not report malformed input:\n%s", output.String())
	}
}

// TestRunStopsAtEndOfInput checks EOF ends the game loop without repeated prompts.
func TestRunStopsAtEndOfInput(t *testing.T) {
	cli := NewCLI(game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O')))
	cli.input = bufio.NewReader(strings.NewReader(""))
	var output bytes.Buffer
	cli.output = &output

	cli.Run()

	if got := strings.Count(output.String(), "Enter row and column:"); got != 1 {
		t.Errorf("prompt count = %d at EOF; want 1", got)
	}
	if strings.Contains(output.String(), "Invalid input:") {
		t.Errorf("Run() reported EOF as invalid input:\n%s", output.String())
	}
}

// TestRunExitsWhenPlayerEntersQ checks that q exits without an invalid-input message.
func TestRunExitsWhenPlayerEntersQ(t *testing.T) {
	match := game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O'))
	cli := NewCLI(match)
	cli.input = bufio.NewReader(strings.NewReader("q\n"))
	var output bytes.Buffer
	cli.output = &output

	cli.Run()

	if match.State() != game.Playing {
		t.Errorf("State() = %v after quitting; want Playing", match.State())
	}
	if strings.Contains(output.String(), "Invalid input:") {
		t.Errorf("Run() treated quit as invalid input:\n%s", output.String())
	}
	if got := strings.Count(output.String(), "Enter row and column:"); got != 1 {
		t.Errorf("prompt count = %d after quitting; want 1", got)
	}
	if got := strings.Count(output.String(), "Enter q + Enter at any time to quit."); got != 1 {
		t.Errorf("quit reminder count = %d; want 1", got)
	}
}

// TestRunCanRetryAfterWinOrDraw checks that either result can start a fresh game.
func TestRunCanRetryAfterWinOrDraw(t *testing.T) {
	tests := []struct {
		name   string
		moves  string
		result string
	}{
		{
			name:   "win",
			moves:  "0 0\n1 0\n0 1\n1 1\n0 2\n",
			result: "Player Sandesh (X) wins!",
		},
		{
			name:   "draw",
			moves:  "0 0\n1 1\n2 2\n0 1\n0 2\n2 0\n1 0\n1 2\n2 1\n",
			result: "It's a draw!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := game.NewGame(game.NewPlayer("Sandesh", 'X'), game.NewPlayer("Rahul", 'O'))
			cli := NewCLI(match)
			cli.input = bufio.NewReader(strings.NewReader(test.moves + "r\nq\n"))
			var output bytes.Buffer
			cli.output = &output

			cli.Run()

			if !strings.Contains(output.String(), test.result) {
				t.Errorf("Run() output missing terminal result %q:\n%s", test.result, output.String())
			}
			if !strings.Contains(output.String(), "Press r + Enter to retry or q + Enter to quit:") {
				t.Errorf("Run() did not prompt for retry or quit:\n%s", output.String())
			}
			if match.State() != game.Playing {
				t.Errorf("State() after retry and quit = %v; want Playing", match.State())
			}
			if match.CurrentPlayer().Symbol != 'X' {
				t.Errorf("current player after retry = %q; want X", match.CurrentPlayer().Symbol)
			}
			if cell, err := match.CellAt(0, 0); err != nil {
				t.Fatalf("CellAt() returned an error: %v", err)
			} else if cell != 0 {
				t.Errorf("board cell after retry = %q; want empty", cell)
			}
		})
	}
}
