package main

import (
	"fmt"
	"io"
	"math/rand"
)

const (
	/* ChessMaxX controls most of the bounding in the
	 * X direction for this program
	 */
	ChessMaxX = 8

	/* ChessMaxY controls most of the bounding in the
	 * Y direction for this program
	 */
	ChessMaxY = 8

	// QUEEN is a constant to represent a queen on a chessboard
	QUEEN     = 'Q'

	// EMPTY is a constant to represent an empty space on a chessbaord
	EMPTY     = ' '
)

/* struct chessboard is the data container that
 * holds all chessboard related information including
 * the data for the boards visualization
 * a list of all queens on the board
 * and the fitness of the board as reported by
 * board.numAttacks
 */
type chessboard struct {
	board [ChessMaxX][ChessMaxY]byte
	queens []*queen
	fitness int
}

/* newChessboard returns a empty chessboard instance
 * @returns an empty chessboard
 */
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

/* newRandomChessboard creates a chessboard with ChessMaxX queens 
 * in unique columns. It guarantees that only ChessMaxX queens will
 * be generated and placed on the board, each in their own unique X column
 * @returns a chessboard full of 8 queens in random positions
 */
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

/* isInBounds is a utility function to determin if a queen
 * is valid given our global constances of ChessMaxX and ChessMaxY
 * @returns an error when the queen is out of bounds
 */
func isInBounds(queen *queen) error {
	if queen.x > ChessMaxX || queen.y > ChessMaxX || queen.x < 1 || queen.y < 1 {
		return newInvalidBoardLocationError(queen.x, queen.y)
	}

	return nil
}

/* isTaken returns and error if the queens passed is at a location
 * that is already taken on the board
 * @param queen the queen to be checked
 * @returns an error dictating if the position is taken, or if the queen passed is invalid
 */
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

/* isEmpty returns an error if the passed queen is at a location that is empty
 * @param queen the queen that needs to be checked against
 * @returns an error dictating if the position is empty or if the queen passed is invalid
 */
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

/* placeQueen places a given queen on the given chessboard
 * @param queen the queen to be placed on the chessboard
 * @returns error if the passed queen is invalid
 */
func (board *chessboard) placeQueen(queen *queen) (err error) {
	err = board.isTaken(queen)
	if err != nil {
		return
	}

	board.board[queen.x-1][queen.y-1] = QUEEN
	board.queens = append(board.queens, queen)
	return
}

/* removeQueen removes a queen on the chessboard represented by the
 * passed queen if it exists
 * @param queen the queen to be removed from the chessboard
 * @returns error if the passed queen is invalid
 */
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

/* numAttacks returns the number of queens that would attack each other if all
 * queens began attacking each other.
 * @returns the number of queens that would attack each other if all queens started attacking
 */
func (board *chessboard) numAttacks() int {
	numAttacks := 0

	for _, queen := range board.queens {
		numAttacks += len(board.attack(queen))
	}

	return numAttacks
}

/* attack returns the list of queens that would be attacked by the passed queen
 * if the passed queen were to attack on the board
 * @param attackingQueen the proposed queen to attack
 * @returns a list of queens that were attacked by the attackingQueen
 */
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

/* print writes out the chessboard to the given io.Writer
 * @param output the output writer to write to
 */
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
