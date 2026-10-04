# Tic-Tac-Toe

A two-player Tic-Tac-Toe game implemented in Go. Player 1 uses `X` and Player 2 uses `O`.

## Requirements

- Go 1.25 or later

## Run

From the project root, start the game with:

```sh
go run ./cmd
```

## How to play

- Enter a row and column as two integers from `0` to `2`, separated by a space. For example, `1 2` selects the middle row and right column.
- Moves outside the board or in an occupied cell are rejected; the same player can try again.
- Enter `q` and press Enter at any time to quit.
- After a win or draw, enter `r` and press Enter to start a new game, or `q` and press Enter to quit.

## Run tests

From the project root:

```sh
go test ./...
```