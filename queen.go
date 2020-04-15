package main

type queen struct {
	x int
	y int
}

func newQueen(x, y int) *queen {
	return &queen{x, y}
}
