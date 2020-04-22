package main

/* 
 * struct queen is a data container
 * representing the logical numebrs
 * of where the queen is located at
 * on the chessboard
 */
type queen struct {
	x int
	y int
}

/* 
 * newQueen is a constructor for queen
 * @param x the cartesian x coordinate for the queen
 * @param y the cartesian y coordiante for the queen
 * @return a pointer to the initialized queen struct
 */
func newQueen(x, y int) *queen {
	return &queen{x, y}
}
