package main

import "fmt"

const (
	ChessMaxX = 8
	ChessMaxY = 8
	QUEEN     = 'Q'
	EMPTY     = ' '
)

type chessboard struct {
	board [ChessMaxX][ChessMaxY]byte
}

func newChessboard() *chessboard {
	board := &chessboard{}
	for row := 0; row < ChessMaxY; row++ {
		for index := range board.board[row] {
			board.board[row][index] = EMPTY
		}
	}

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
	return
}

func (board *chessboard) removeQueen(queen *queen) (err error) {
	err = board.isEmpty(queen)
	if err != nil {
		return
	}

	board.board[queen.x-1][queen.y-1] = EMPTY
	return
}

func (board *chessboard) attacks(attackingQueen *queen) (attackedQueens []*queen, err error) {
	attackedQueens = make([]*queen, 0)

	return
}

func (board *chessboard) print() {
	const (
		whiteBackground = "\033[48;5;15m"
		whiteForeground = "\033[38;5;15m"
		blackBackground = "\033[48;5;0m"
		blackForeground = "\033[38;5;0m"
	)

	currentBackground := whiteBackground
	currentForeground := blackForeground

	for row := 0; row < ChessMaxY; row++ {
		fmt.Printf("%d|", row+1)
		fmt.Print("\033[1m")
		for _, element := range board.board[row] {
			fmt.Printf("%s ", currentBackground)
			fmt.Printf("%s%c", currentForeground, element)
			fmt.Printf("%s ", currentBackground)

			if currentBackground == whiteBackground {
				currentBackground = blackBackground
				currentForeground = whiteForeground
			} else {
				currentBackground = whiteBackground
				currentForeground = blackForeground
			}
		}
		if currentBackground == whiteBackground {
			currentBackground = blackBackground
			currentForeground = whiteForeground
		} else {
			currentBackground = whiteBackground
			currentForeground = blackForeground
		}

		fmt.Println("\033[m|")
	}
	fmt.Print("  ")
	for i := 1; i < 9; i++ {
		fmt.Printf(" %d ", i)
	}
	fmt.Println()
}
