package main

import "os"

func main() {
	logfile, err := os.OpenFile("8queens.log", os.O_TRUNC|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()

	EnableLogging(logfile)

	chessboard := newChessboard()
	queen := newQueen(1, 6)
	if err := chessboard.placeQueen(queen); err != nil {
		panic(err)
	}
	chessboard.print()
}
