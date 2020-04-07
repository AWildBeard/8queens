package main

import "fmt"

type boardLocationTakenError struct {
	y int
	x int
}

func newBoardLocationTakenError(x, y int) *boardLocationTakenError {
	return &boardLocationTakenError{x, y}
}

func (blte *boardLocationTakenError) Error() string {
	return fmt.Sprintf("location %d, %d is already taken", blte.x, blte.y)
}

type boardLocationEmptyError struct {
	y int
	x int
}

func newBoardLocationEmptyError(x, y int) *boardLocationEmptyError {
	return &boardLocationEmptyError{x, y}
}

func (blee *boardLocationEmptyError) Error() string {
	return fmt.Sprintf("location %d, %d is already empty", blee.x, blee.y)
}

type invalidBoardLocationError struct {
	y int
	x int
}

func newInvalidBoardLocationError(x, y int) *invalidBoardLocationError {
	return &invalidBoardLocationError{x, y}
}

func (ible *invalidBoardLocationError) Error() string {
	return fmt.Sprintf("location %d, %d is invalid", ible.x, ible.y)
}
