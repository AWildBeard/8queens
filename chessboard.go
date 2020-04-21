package main

import (
	"fmt"
	"io"
	"math/rand"
)

const (
	ChessMaxX = 8
	ChessMaxY = 8
	QUEEN     = 'Q'
	EMPTY     = ' '
)

type chessboard struct {
	board [ChessMaxX][ChessMaxY]byte
	queens []*queen
	fitness int
}

func newChessboard() *chessboard {
	board := &chessboard{}
	for row := 0; row < ChessMaxY; row++ {
		for index := range board.board[row] {
			board.board[row][index] = EMPTY
		}
	}

	board.queens = make([]*queen, 0)

	return board
}

// newRandomChessboard creates a chessboard with ChessMaxX queens 
// in unique columns. It guarantees that only ChessMaxX queens
func newRandomChessboard() *chessboard {
	board := &chessboard{}
	board.queens = make([]*queen, 0)

	for column := 0 ; column < ChessMaxX ; column++ {
		randomRowLocation := rand.Intn(ChessMaxY)
		for row := 0 ; row < ChessMaxY ; row++ {
			if row == randomRowLocation {
				board.board[column][row] = QUEEN
				board.queens = append(board.queens, newQueen(column+1, row+1))
			} else {
				board.board[column][row] = EMPTY
			}
		}
	}

	board.fitness = board.numAttacks()

	return board
}

func isInBounds(queen *queen) error {
	if queen.x > ChessMaxX || queen.y > ChessMaxX || queen.x < 1 || queen.y < 1 {
		return newInvalidBoardLocationError(queen.x, queen.y)
	}

	return nil
}

func (board *chessboard) isTaken(queen *queen) (err error) {
	err = isInBounds(queen)
	if err != nil {
		return
	}

	if board.board[queen.x-1][queen.y-1] == QUEEN {
		err = newBoardLocationTakenError(queen.x, queen.y)
	}

	return
}

func (board *chessboard) isEmpty(queen *queen) (err error) {
	err = isInBounds(queen)
	if err != nil {
		return
	}

	if board.board[queen.x-1][queen.y-1] == EMPTY {
		err = newBoardLocationEmptyError(queen.x, queen.y)
	}

	return
}

func (board *chessboard) placeQueen(queen *queen) (err error) {
	err = board.isTaken(queen)
	if err != nil {
		return
	}

	board.board[queen.x-1][queen.y-1] = QUEEN
	board.queens = append(board.queens, queen)
	return
}

func (board *chessboard) removeQueen(queen *queen) (err error) {
	err = board.isEmpty(queen)
	if err != nil {
		return
	}

	board.board[queen.x-1][queen.y-1] = EMPTY
	for index, boardQueen := range board.queens {
		if *queen == *boardQueen {
			board.queens[index] = board.queens[len(board.queens) - 1]
			board.queens = board.queens[:len(board.queens) - 1]
		}
	}
	return
}

func (board *chessboard) numAttacks() int {
	numAttacks := 0

	for _, queen := range board.queens {
		numAttacks += len(board.attack(queen))
	}

	return numAttacks
}

func (board *chessboard) attack(attackingQueen *queen) (attackedQueens []*queen) {
	attackedQueens = make([]*queen, 0)

	checkAndAddQueen := func(x, y int) bool {
		candidateQueen := newQueen(x, y)

		// assert that the reason the attack failed, is because there is a queen
		err := board.isTaken(candidateQueen)
		if _, ok := err.(*boardLocationTakenError); ok {
			attackedQueens = append(attackedQueens, candidateQueen)
			return true
		} else if err != nil {
			panic(err) // Panic if bounds are exceeded/ unknown error
		}

		return false
	}

	// Check horizontal left
	for xIndex := attackingQueen.x - 1; /* skip the attacker */ xIndex > 0; xIndex-- {
		if checkAndAddQueen(xIndex, attackingQueen.y) {
			break
		}
	}

	// Check horizontal right
	for xIndex := attackingQueen.x + 1; /* skip the attacker */ xIndex <= ChessMaxX; xIndex++ {
		if checkAndAddQueen(xIndex, attackingQueen.y) {
			break
		}
	}

	// Check vertical down
	for yIndex := attackingQueen.y - 1; /* skip the attacker */ yIndex > 0; yIndex-- {
		if checkAndAddQueen(attackingQueen.x, yIndex) {
			break
		}
	}

	// Check vertical up
	for yIndex := attackingQueen.y + 1; /* skip the attacker */ yIndex <= ChessMaxY; yIndex++ {
		if checkAndAddQueen(attackingQueen.x, yIndex) {
			break
		}
	}

	// Check Diagonal to Top Left
	for xIndex, yIndex := attackingQueen.x-1, attackingQueen.y+1; xIndex > 0 && yIndex <= ChessMaxY; xIndex, yIndex = xIndex-1, yIndex+1 {
		if checkAndAddQueen(xIndex, yIndex) {
			break
		}
	}

	// Check Diagonal to Bottom Right
	for xIndex, yIndex := attackingQueen.x+1, attackingQueen.y-1; xIndex <= ChessMaxX && yIndex > 0; xIndex, yIndex = xIndex+1, yIndex-1 {
		if checkAndAddQueen(xIndex, yIndex) {
			break
		}
	}

	// Check Diagonal to Top Right
	for xIndex, yIndex := attackingQueen.x+1, attackingQueen.y+1; xIndex <= ChessMaxX && yIndex <= ChessMaxY; xIndex, yIndex = xIndex+1, yIndex+1 {
		if checkAndAddQueen(xIndex, yIndex) {
			break
		}
	}

	// Check Diagonal to Bottom Left
	for xIndex, yIndex := attackingQueen.x-1, attackingQueen.y-1; xIndex > 0 && yIndex > 0; xIndex, yIndex = xIndex-1, yIndex-1 {
		if checkAndAddQueen(xIndex, yIndex) {
			break
		}
	}

	return
}

func (board *chessboard) print(output io.Writer) {
	const (
		whiteBackground = "\033[48;5;15m"
		whiteForeground = "\033[38;5;15m"
		blackBackground = "\033[48;5;0m"
		blackForeground = "\033[38;5;0m"
	)

	currentBackground := whiteBackground
	currentForeground := blackForeground

	rotateColors := func() {
		if currentBackground == whiteBackground {
			currentBackground = blackBackground
			currentForeground = whiteForeground
		} else {
			currentBackground = whiteBackground
			currentForeground = blackForeground
		}
	}

	fmt.Fprintln(output)
	for column := ChessMaxX - 1; column >= 0; column-- {
		fmt.Fprintf(output, "%d|", column+1)
		fmt.Fprint(output, "\033[1m")
		for row := 0; row < ChessMaxY; row++ {
			fmt.Fprintf(output, "%s ", currentBackground)
			fmt.Fprintf(output, "%s%c", currentForeground, board.board[row][column])
			fmt.Fprintf(output, "%s ", currentBackground)

			rotateColors()
		}
		fmt.Fprintln(output, "\033[m|")

		rotateColors()
	}

	fmt.Fprint(output, "  ")
	for i := 1; i < 9; i++ {
		fmt.Fprintf(output, " %d ", i)
	}

	fmt.Fprintln(output)
}
