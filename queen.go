package main

type queen struct {
	y int
	x int
}

func newQueen(x, y int) *queen {
	return &queen{x, y}
}
